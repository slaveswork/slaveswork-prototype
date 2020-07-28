package main

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