package main

import "fmt"

type Tile struct {
	Index int `json:"index"`
	Xmin  int `json:"xmin"`
	Ymin  int `json:"ymin"`
	Xmax  int `json:"xmax"`
	Ymax  int `json:"ymax"`
	Frame int `json:"fram"`
}

func (t *Tile) prettyPrint(){
	fmt.Printf("Tile: { \"index\" : %d, \"xmin\" : %d, \"ymin\" : %d, \"xmax\" : %d, \"ymax\" : %d, \"frame\" : %d}\n", t.Index, t.Xmin, t.Ymin, t.Xmax, t.Ymax, t.Frame)
}

