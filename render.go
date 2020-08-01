package main

import (
	"log"
)

func (h *Host) DoRender() {
	path := <-h.filePath // blend file path

	var successful int
	var result bool

	for {
		for i, _ := range h.tiles {
			if !h.tiles[i].Active {
				h.tiles[i].Active = true
				result = h.tiles[i].Dispatch(h, path)

				if result == false {
					h.tiles[i].Active = false
					log.Println("Dispatch failed Tile index :", h.tiles[i].Index)
					<- h.freeWorker // Wait for free worker channel
				}
			}
		}

		successful = 0
		for _, tile := range h.tiles {
			if tile.Success {
				successful += 1
			}
		}

		// end code
		if (successful / len(h.tiles)) == 1.0 {
			break
		}
	}
}
