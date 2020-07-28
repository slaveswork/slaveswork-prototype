package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type Tile struct {
	Index string `json:"index"`
	Xmin  string `json:"xmin"`
	Ymin  string `json:"ymin"`
	Xmax  string `json:"xmax"`
	Ymax  string `json:"ymax"`
	Frame string `json:"fram"`

	// In JSON parsing, the following parameters are ignored.
	Active  bool `json:"-"`
	Success bool `json:"-"`
}

func (t *Tile) prettyPrint() {
	fmt.Printf("Tile: { \"index\" : %s, \"xmin\" : %s, \"ymin\" : %s, \"xmax\" : %s, \"ymax\" : %s, \"frame\" : %s}\n", t.Index, t.Xmin, t.Ymin, t.Xmax, t.Ymax, t.Frame)
}

func (t *Tile) Dispatch(h *Host, path string) {
	var worker *Worker

	for key := range h.workers {
		worker = h.workers[key]
		if worker.Status == "" {
			r, w := io.Pipe()
			m := multipart.NewWriter(w)

			go func() {
				defer w.Close()
				defer m.Close()

				part, err := m.CreateFormFile("blend-file", "main.blend")

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
			req.Header.Add("index", t.Index)
			req.Header.Add("xmin", t.Xmin)
			req.Header.Add("ymin", t.Ymin)
			req.Header.Add("xmax", t.Xmax)
			req.Header.Add("ymax", t.Ymax)
			req.Header.Add("fram", t.Frame)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal("func : Dispatch\n", err)
			}

			fmt.Println("worker ", worker.Id, " dispatch Status : ",resp.Status)

			resp.Body.Close()
			worker.Status = "running"
			break
		}
	}
}
