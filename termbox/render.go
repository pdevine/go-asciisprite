package termbox

// Diff renderer: emit the minimal escape stream turning the front
// buffer into the back buffer.  Escape dialect:
//
//	CSI {row};{col}H   cursor position (1-based)
//	CSI 0;m            SGR reset
//	CSI 38;5;{n}m      256-color fg, CSI 48;5;{n}m bg
//	CSI 39m / 49m      default fg/bg
//	SGR 1 / 4 / 7      bold, underline, reverse
//	CSI 30-37/40-47 m  16-color fg/bg (OutputNormal mode)

import (
	"io"
	"strconv"
	"unicode/utf8"
)

type renderer struct {
	out     io.Writer
	outMode OutputMode

	cx, cy  int // cursor position of the emitted stream, -1 = unknown
	fg, bg  Attribute
	bold    bool
	underln bool
	reverse bool
	haveSGR bool // fg/bg/bits fields are valid
	scrW    int
	scrH    int
}

func newRenderer(w io.Writer) *renderer {
	return &renderer{out: w, cx: -1, cy: -1, outMode: Output256}
}

func (r *renderer) write(s string) { io.WriteString(r.out, s) }

func (r *renderer) cup(x, y int) {
	if r.cx == x && r.cy == y {
		return
	}
	r.write("\x1b[" + strconv.Itoa(y+1) + ";" + strconv.Itoa(x+1) + "H")
	r.cx, r.cy = x, y
}

func (r *renderer) sgrReset() {
	r.write("\x1b[0m")
	r.fg, r.bg = ColorDefault, ColorDefault
	r.bold, r.underln, r.reverse = false, false, false
	r.haveSGR = true
}

// colorIndex maps an Attribute color field to a palette index 0-255,
// or -1 for default.
func colorIndex(a Attribute) int {
	c := int(a) & 0x1ff
	if c == 0 {
		return -1
	}
	// the named 16-color constants occupy slots 1-16 (0 is default) but
	// map to palette indices 0-15
	return c - 1
}

func (r *renderer) applyColor(a Attribute, foreground bool) {
	idx := colorIndex(a)
	if r.outMode == OutputNormal && idx >= 16 {
		idx %= 16
	}
	if idx < 0 {
		if foreground {
			r.write("\x1b[39m")
		} else {
			r.write("\x1b[49m")
		}
		return
	}
	if r.outMode == OutputNormal {
		base := 40
		if foreground {
			base = 30
		}
		r.write("\x1b[" + strconv.Itoa(base+idx) + "m")
		return
	}
	pfx := "48"
	if foreground {
		pfx = "38"
	}
	r.write("\x1b[" + pfx + ";5;" + strconv.Itoa(idx) + "m")
}

func (r *renderer) setStyle(fg, bg Attribute) {
	wantBold := (fg|bg)&AttrBold != 0
	wantUnd := (fg|bg)&AttrUnderline != 0
	wantRev := (fg|bg)&AttrReverse != 0
	fgc, bgc := fg&0x1ff, bg&0x1ff

	if !r.haveSGR ||
		(r.bold && !wantBold) || (r.underln && !wantUnd) || (r.reverse && !wantRev) {
		// dropping an attr bit requires a reset first
		r.sgrReset()
	}
	curFg, curBg := r.fg&0x1ff, r.bg&0x1ff
	if curFg != fgc {
		r.applyColor(fg, true)
		r.fg = fg
	}
	if curBg != bgc {
		r.applyColor(bg, false)
		r.bg = bg
	}
	if wantBold && !r.bold {
		r.write("\x1b[1m")
		r.bold = true
	}
	if wantUnd && !r.underln {
		r.write("\x1b[4m")
		r.underln = true
	}
	if wantRev && !r.reverse {
		r.write("\x1b[7m")
		r.reverse = true
	}
}

// render emits the escape stream converting front into back and swaps.
// Returns the bytes written count, for tests/diagnostics.
func (r *renderer) render(front, back *screenBuf) int {
	before := 0
	if cw, ok := r.out.(interface{ Len() int }); ok {
		before = cw.Len()
	}
	for y := 0; y < back.height; y++ {
		yw := y * back.width
		for x := 0; x < back.width; x++ {
			c := back.cells[yw+x]
			if front.width == back.width && front.height == back.height &&
				front.cells[yw+x] == c {
				continue
			}
			r.cup(x, y)
			r.setStyle(c.fg, c.bg)
			var utf [8]byte
			n := utf8.EncodeRune(utf[:], c.ch)
			r.out.Write(utf[:n])
			r.cx = x + 1 // cursor advances one cell after the write
		}
	}
	// swap: back becomes the new front; old front storage becomes back
	front.cells, back.cells = back.cells, front.cells
	front.width, front.height = back.width, back.height
	if cw, ok := r.out.(interface{ Len() int }); ok {
		return cw.Len() - before
	}
	return -1
}
