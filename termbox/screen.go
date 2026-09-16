package termbox

// Cell grid: front (on screen) and back (app-written) buffers, plus
// dirty tracking for the diff renderer.

type cell struct {
	ch     rune
	fg, bg Attribute
}

type screenBuf struct {
	cells  []cell
	width  int
	height int
}

func newScreenBuf(w, h int) screenBuf {
	b := screenBuf{width: w, height: h, cells: make([]cell, w*h)}
	b.clear()
	return b
}

func (b *screenBuf) clear() {
	for i := range b.cells {
		b.cells[i] = cell{ch: ' ', fg: ColorDefault, bg: ColorDefault}
	}
}
