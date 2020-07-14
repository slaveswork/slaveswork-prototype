package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type Tile struct {
	Index int `json:"index"`
	Xmin  int `json:"xmin"`
	Ymin  int `json:"ymin"`
	Xmax  int `json:"xmax"`
	Ymax  int `json:"ymax"`
	Frame int `json:"fram"`

	// In JSON parsing, the following parameters are ignored.
	Active  bool `json:"-"`
	Success bool `json:"-"`
}

func (t *Tile) prettyPrint() {
	fmt.Printf("Tile: { \"index\" : %d, \"xmin\" : %d, \"ymin\" : %d, \"xmax\" : %d, \"ymax\" : %d, \"frame\" : %d}\n", t.Index, t.Xmin, t.Ymin, t.Xmax, t.Ymax, t.Frame)
}

func (t *Tile) Dispatch(h *Host, path string) {
	for _, worker := range h.workers {
		if worker.Status == "" {
			r, w := io.Pipe()
			m := multipart.NewWriter(w)

			go func() {
				defer w.Close()
				defer m.Close()

				part, err := m.CreateFormFile("blend-file", path)
				if err != nil {
					log.Fatal("func : Dispatch\n", err)
				}

				file, err := os.Open(path)
				if err != nil {
					log.Fatal("func : Dispatch\n", err)
				}
				defer file.Close()

				if _, err := io.Copy(part, file); err != nil {
					log.Fatal("func : Dispatch\n", err)
				}
			}()

			req, err := http.NewRequest("POST", worker.Address.generateHostAddress("/render/resource"), r)
			if err != nil {
				log.Fatal("func : Dispatch\n", err)
			}

			req.Header.Add("Content-Type", m.FormDataContentType())
			req.Header.Add("index", strconv.Itoa(t.Index))
			req.Header.Add("xmin", strconv.Itoa(t.Xmin))
			req.Header.Add("ymin", strconv.Itoa(t.Ymin))
			req.Header.Add("xmax", strconv.Itoa(t.Xmax))
			req.Header.Add("ymax", strconv.Itoa(t.Ymax))
			req.Header.Add("fram", strconv.Itoa(t.Frame))

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal("func : Dispatch\n", err)
			}

			resp.Body.Close()
			worker.Status = "running"
			break
		}
	}
}

func (t *Tile) makePythonFile() string {

}