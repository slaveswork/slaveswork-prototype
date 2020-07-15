package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (h *Host) ReceiveTaskResource(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("func : ReceiveTaskResource\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	fmt.Println("blen : ", string(reqBody[0:4]))
	if blen := string(reqBody[:4]); strings.Compare(blen, "blen") != 0 {
		log.Fatal("func : ReceiveTaskResource\n", "This is not a task resource request.")
	}
	buffer := bytes.NewBuffer(reqBody[4:8])
	fileSize, err := binary.ReadVarint(buffer)
	fmt.Println("size : ", fileSize)
	if err != nil {
		log.Fatal("func : ReceiveTaskResource\n", "This is not a number.")
	}

	//if fileSize != len(reqBody[8:]) {
	//	log.Fatal("func : ReceiveTaskResource\n", "An error occurred while receiving the file.")
	//}

	dir := filepath.Join(os.TempDir(), "main.blend")
	err = ioutil.WriteFile(dir, reqBody[8:], 0644)
	if err != nil {
		log.Fatal("func : ReceiveTaskResource\n", "File Writing Error.")
	}
	fmt.Println("dir:", dir)
	h.filePath <- dir
}

func (h *Host) DoRender() {
	path := <- h.filePath // blend file path

	var remaining []Tile
	for {
		// remaining tiles -> should be render.
		remaining = filter(h.tiles, func(v Tile) bool {
			return !v.Active
		})

		// When there are no tiles left, rendering complete
		if len(remaining) == 0 {
			break
		}

		// for all remaining tiles
		for _, tile := range remaining {
			tile.Active = true
			tile.Dispatch(h, path)
		}
	}
}

func filter(vs []Tile, f func(Tile) bool) []Tile {
	vsf := make([]Tile, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}