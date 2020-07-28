package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/Equanox/gotron"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Worker struct {
	window      *gotron.BrowserWindow `json:"-"` // to ignore in marshaling
	hostAddress Address               `json:"-"` // to ignore in marshaling
	ci          chan int              `json:"-"`

	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Method  string  `json:"method"`
	Address Address `json:"address"`

	// Worker status -> ignore
	Status      string      `json:"-"`
	T           chan Tile   `json:"-"`
	FilePath    chan string `json:"-"`
	Config      Config      `json:"-"`
	BlenderPath string      `json:"-"`
}

func newWorker(w *gotron.BrowserWindow) *Worker {
	return &Worker{
		window:   w, // for worker's window.
		ci:       make(chan int),
		T:        make(chan Tile),
		FilePath: make(chan string),
	}
}

func (w *Worker) run() {
	listener := w.init()
	w.gotronMessageHandler()
	w.sendConfig()
	go w.sendWorkerStatus()
	go w.httpMessageHandler(listener)
	go w.renderTileWithBlender()
}

func (w *Worker) init() (listener net.Listener) {
	w.Name, _ = os.Hostname()
	w.Address, listener = newAddress() // initialize worker's address.
	w.Config = UnmarshalConfig()
	return
}

func (w *Worker) gotronMessageHandler() {
	w.window.On(&gotron.Event{Event: "app.connect.device"}, w.sendConnectionRequest) // connection between worker and host.
	w.window.On(&gotron.Event{Event: "app.blender.path"}, w.receiveBlenderPath)      // blender.exe path
}

func (w *Worker) httpMessageHandler(listener net.Listener) {
	http.HandleFunc("/render/resource", w.receiveRenderResource)

	http.Serve(listener, nil)
}

func (w *Worker) sendConnectionRequest(bin []byte) {
	var receivedMessage GotronMessage // for Unmarshal window message
	body := struct {
		Address
		Token string `json:token` // Token...
	}{}
	receivedMessage.Body = &body

	// "app.connect.device" Unmarshal JSON message
	if err := json.Unmarshal(bin, &receivedMessage); err != nil {
		log.Fatal("func : sendConnectionRequest\n", err)
	}
	w.hostAddress = body.Address // Initialize host address.
	w.Method = "add"             // first send "add" worker to host.

	w.Config.HostIp = body.Address.IP
	w.Config.SaveConfig()

	// Marshaling body for connection request.
	requestBody, err := json.MarshalIndent(w, "", "    ")
	if err != nil {
		log.Fatal("func : sendConnectionRequest\n", err)
	}

	// send request to host.
	url := body.Address.generateHostAddress(body.Token) // make host url --> format : http://{ip}:{port}/{token}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal("func : sendConnectionRequest\n", err)
	}

	// resp's body will have this worker's Id(for Host's worker management).
	defer resp.Body.Close()

	respId := struct { // temporary struct for receiving response body.
		Id int `json:"id"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&respId); err != nil {
		log.Fatal("func : sendConnectionRequest\n", err)
	}

	w.Id = respId.Id // set up worker's Id for host's management.
	w.ci <- respId.Id
}

func (w *Worker) sendWorkerStatus() {
	<-w.ci      // wait function until setting worker's Id.
	close(w.ci) // close channel... we don't need that anymore.

	url := w.hostAddress.generateHostAddress("status") // Host address.
	w.Method = "update"                                // from now on, send "update" method to host.

	for range time.Tick(time.Second * 5) { // repeat sending this worker's status every 5 seconds.
		requestBody, err := json.MarshalIndent(w, "", "    ")
		if err != nil {
			log.Fatal("func : sendWorkerStatus\n", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Fatal("func : sendWorkerStatus\n", err)
		}

		resp.Body.Close() // close response body.
	}
}

func (w *Worker) sendConfig() {
	message := GotronMessage{
		Event: &gotron.Event{Event: "window.send.config"},
	}
	message.Body = struct { // temporary struct for sending token message.
		HostIp      string `json:"hostIp"`
		BlenderPath string `json:"blenderPath"`
	}{
		HostIp:      w.Config.HostIp,
		BlenderPath: w.Config.BlenderPath, // initialize value.
	}

	//checkJSON(message) // Printing Message for validation.
	w.window.Send(message)
}

func (w *Worker) receiveBlenderPath(bin []byte) {
	var message GotronMessage
	body := struct {
		BlenderPath string `json:blenderPath`
	}{}
	message.Body = &body

	if err := json.Unmarshal(bin, &message); err != nil {
		log.Fatal("func : receiveBlenderPath\n", err)
	}
	w.Config.BlenderPath = body.BlenderPath
	w.Config.SaveConfig()
}

func (w *Worker) receiveRenderResource(rw http.ResponseWriter, r *http.Request) {
	t := Tile{
		Index: r.Header.Get("index"),
		Xmin:  r.Header.Get("xmin"),
		Ymin:  r.Header.Get("ymin"),
		Xmax:  r.Header.Get("xmax"),
		Ymax:  r.Header.Get("ymax"),
		Frame: r.Header.Get("fram"),
	}

	file, header, err := r.FormFile("blend-file")
	if err != nil {
		log.Fatal("func : receiveRenderResource\n", err)
	}

	defer file.Close()

	path := filepath.Join(os.TempDir(), header.Filename)
	out, err := os.Create(path)
	if err != nil {
		log.Fatal(rw, "Unable to create the file for writing.")
		return
	}

	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(rw, err)
	}

	w.FilePath <- path
	w.T <- t
}

func (w *Worker) renderTileWithBlender() {
	var outputFile string
	var outputPath string

	var result chan string
	var done string

	for {
		result = make(chan string)
		blendFile := <-w.FilePath // *.blend file path
		tile := <-w.T

		outputFile = "output" + tile.Index + ".exr"
		outputPath = filepath.Join(os.TempDir(), outputFile)

		cmd := exec.Command(w.Config.BlenderPath,
			"-b", blendFile,
			"-F", "EXR",
			"--render-output", outputPath,
			"-Y",
			"-noaudio",
			"-E", "CYCLES",
			"-P", "worker.py",
			"--", tile.Index, tile.Xmin, tile.Ymin, tile.Xmax, tile.Ymax, tile.Frame)

		stdout, err := cmd.StdoutPipe()
		stderr, err := cmd.StderrPipe()
		err = cmd.Start()
		if err != nil {
			log.Fatal("func : renderTileWithBlender\n", err)
		}

		go w.copyBlenderOutput(stdout, result)
		go w.copyBlenderOutput(stderr, result)
		cmd.Wait()

		done = <-result

		req, _ := http.NewRequest("POST", w.hostAddress.generateHostAddress("/task/result"), nil)

		if _, err := os.Stat(outputPath); err == nil {
			pr, pw := io.Pipe()
			m := multipart.NewWriter(pw)

			go func() {
				defer pw.Close()
				defer m.Close()

				part, err := m.CreateFormFile("blend-file", outputFile)
				file, err := os.Open(outputPath)
				if err != nil {
					log.Fatal("func : renderTileWithBlender\n", err)
				}

				defer file.Close()

				if _, err := io.Copy(part, file); err != nil {
					log.Fatal("func : renderTileWithBlender\n", err)
				}
			}()

			req.Body = pr
		}

		req.Header.Add("Worker", strconv.Itoa(w.Id))
		req.Header.Add("Tile", tile.Index)
		req.Header.Add("Result", done)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		resp.Body.Close()
		close(result)
	}
}

func (w *Worker) copyBlenderOutput(r io.Reader, ch chan string) {
	sc := bufio.NewScanner(r)

	var text string

	for sc.Scan() {
		text = sc.Text()

		log.Println(text)

		// When 'sc' get the rendering complete text.
		if strings.Contains(text, "Success.") {
			ch <- "success"
		} else if strings.Contains(text, "Error.") {
			ch <- text // TODO find another way to get error message.
		}
	}
}
