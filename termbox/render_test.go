package termbox

// Golden tests for the diff renderer: exact byte streams for known cell
// contents, written via a bytes.Buffer (render is io.Writer based — no
// tty needed).

import (
	"bytes"
	"testing"
)

func renderOnce(b *screenBuf) string {
	front := newScreenBuf(b.width, b.height)
	buf := &bytes.Buffer{}
	r := newRenderer(buf)
	r.render(&front, b)
	return buf.String()
}

func TestRenderSingleCell(t *testing.T) {
	b := newScreenBuf(4, 2)
	b.cells[1*4+2] = cell{ch: 'X', fg: ColorRed, bg: ColorDefault}
	got := renderOnce(&b)
	// cursor to (2,1) = CSI 2;3H ; fg red = palette index 9
	want := "\x1b[2;3H\x1b[0m\x1b[38;5;9mX"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRenderCoalescesRuns(t *testing.T) {
	b := newScreenBuf(3, 1)
	for i := range b.cells {
		b.cells[i] = cell{ch: '.', fg: ColorBlue, bg: ColorGreen}
	}
	got := renderOnce(&b)
	// green = palette index 2; blue = 12 (0;10;255 in the xterm table)
	want := "\x1b[1;1H\x1b[0m\x1b[38;5;12m\x1b[48;5;2m..."
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRenderSkipsUnchanged(t *testing.T) {
	front := newScreenBuf(2, 1)
	back := newScreenBuf(2, 1)
	front.cells[0] = cell{ch: 'a', fg: 1, bg: 0}
	back.cells[0] = cell{ch: 'a', fg: 1, bg: 0}
	back.cells[1] = cell{ch: 'b', fg: 2, bg: 0}
	buf := &bytes.Buffer{}
	r := newRenderer(buf)
	r.render(&front, &back)
	if buf.String() != "\x1b[1;2H\x1b[0m\x1b[38;5;1mb" {
		t.Fatalf("got %q", buf.String())
	}
	// after the swap, the front buffer is the just-rendered screen; the
	// (untouched) back buffer is the stale one.  An immediate second
	// Flush redraws the whole frame (all cells differ stale vs new), so
	// to test "unchanged in, nothing out" make back identical to front.
	back.cells[0] = front.cells[0]
	back.cells[1] = front.cells[1]
	buf.Reset()
	r.render(&front, &back)
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got %q", buf.String())
	}
}

func TestRenderOutputNormalClamp16(t *testing.T) {
	b := newScreenBuf(1, 1)
	b.cells[0] = cell{ch: 'X', fg: Attribute(1 + 200), bg: ColorDefault}
	buf := &bytes.Buffer{}
	r := newRenderer(buf)
	r.outMode = OutputNormal
	front := newScreenBuf(1, 1)
	r.render(&front, &b)
	// OutputNormal clamps 200 to 200%16=8 -> SGR 38
	want := "\x1b[1;1H\x1b[0m\x1b[38mX"
	if buf.String() != want {
		t.Fatalf("got %q want %q", buf.String(), want)
	}
}

func TestRenderDefaultColorsUseSGR39_49(t *testing.T) {
	b := newScreenBuf(1, 1)
	b.cells[0] = cell{ch: 'X', fg: ColorRed, bg: ColorRed}
	// start from a state with red fg/bg, then render a default cell
	buf := &bytes.Buffer{}
	r := newRenderer(buf)
	r.haveSGR = true
	r.fg, r.bg = ColorRed, ColorRed
	d := newScreenBuf(1, 1)
	d.cells[0] = cell{ch: 'X', fg: ColorDefault, bg: ColorDefault}
	front := newScreenBuf(1, 1)
	r.render(&front, &d)
	if !bytes.Contains(buf.Bytes(), []byte("\x1b[39m")) ||
		!bytes.Contains(buf.Bytes(), []byte("\x1b[49m")) {
		t.Fatalf("missing default-color resets: %q", buf.String())
	}
}
