package sprite

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"sort"
	"strings"

	palette "github.com/pdevine/go-asciisprite/palette"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

// Blocks provides a map of bits to Unicode block character runes.
var Blocks = map[int]rune{
	0:  ' ',
	1:  '▘',
	2:  '▝',
	3:  '▀',
	4:  '▖',
	5:  '▌',
	6:  '▞',
	7:  '▛',
	8:  '▗',
	9:  '▚',
	10: '▐',
	11: '▜',
	12: '▄',
	13: '▙',
	14: '▟',
	15: '█',
}

type Surface struct {
	Blocks [][]rune
	Width  int
	Height int
	Alpha  bool
}

// ColorMap provides an interpolation map of characters to termbox colors.
var ColorMap = map[rune]tm.Attribute{
	'R': tm.ColorRed,
	'b': tm.Attribute(53),
	't': tm.Attribute(180),
	'Y': tm.ColorYellow,
	'N': tm.ColorBlack,
	'B': tm.ColorBlue,
	'o': tm.Attribute(209),
	'O': tm.Attribute(167),
	'w': tm.ColorWhite,
	'X': tm.ColorWhite,
	'g': tm.ColorGreen,
	'G': tm.Attribute(35),
}

// Convert is a convenience function to create a 1 bit Costume from a string
func Convert(s string) Costume {
	sf := NewSurfaceFromString(s, false)
	return sf.ConvertToCostume()
}

// ColorConvert convenience function to create a color Costume from a string
func ColorConvert(s string, bg tm.Attribute) Costume {
	sf := NewSurfaceFromString(s, false)
	return sf.ConvertToColorCostume(bg)
}

// NewSurface creates a Surface
func NewSurface(width, height int, alpha bool) Surface {
	blocks := make([][]rune, height, height)
	for cnt := 0; cnt < height; cnt++ {
		blocks[cnt] = make([]rune, width, width)
	}

	s := Surface{
		Blocks: blocks,
		Width:  width,
		Height: height,
		Alpha:  alpha,
	}
	return s
}

// NewSurfaceFromString creates a Surface which can be converted to a Costume
func NewSurfaceFromString(s string, alpha bool) Surface {
	l := strings.Split(s, "\n")
	maxR := len(l) + len(l)%2

	// all block sprites must be even
	m := make([][]rune, maxR, maxR)

	var maxC int
	for _, r := range l {
		maxC = max(maxC, len(r)+len(r)%2)
	}

	for rcnt, r := range l {
		m[rcnt] = make([]rune, maxC, maxC)
		for ccnt, c := range r {
			if c == ' ' {
				if !alpha {
					m[rcnt][ccnt] = 0
				} else {
					continue
				}
			} else {
				m[rcnt][ccnt] = c
			}
		}
	}

	// make certain we make a row for any added space
	if len(l) < maxR {
		m[maxR-1] = make([]rune, maxC, maxC)
	}
	sf := Surface{
		Blocks: m,
		Width:  maxC,
		Height: maxR,
		Alpha:  alpha,
	}
	return sf
}

// NewSurfaceFromPng returns a Surface from a PNG file
func NewSurfaceFromPng(fn string, alpha bool) Surface {
	f, err := os.Open(fn)
	if err != nil {
		//
	}

	img, err := png.Decode(f)
	if err != nil {
		//
	}
	return NewSurfaceFromImage(img, alpha)
}

// NewSurfaceFromImage returns a Surface from an image.Image
func NewSurfaceFromImage(img image.Image, alpha bool) Surface {
	bnd := img.Bounds()
	maxR := (bnd.Max.Y - bnd.Min.Y) + (bnd.Max.Y-bnd.Min.Y)%2
	maxC := (bnd.Max.X - bnd.Min.X) + (bnd.Max.X-bnd.Min.X)%2

	// all block sprites must be even
	m := make([][]rune, maxR, maxR)

	for y := 0; y < bnd.Max.Y-bnd.Min.Y; y++ {
		m[y] = make([]rune, maxC, maxC)
		for x := 0; x < bnd.Max.X-bnd.Min.X; x++ {
			c := img.At(x+bnd.Min.X, y+bnd.Min.Y)
			r, g, b, a := c.RGBA()
			// we don't properly support the alpha channel, so only draw the pixel if
			// the alpha is set
			if a > 0 {
				// we only support 256 colour mode, so find the closest colour
				// in the palette and create an entry in our colour map if needed
				i := palette.NearestIndex(color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
				if i > -1 {
					m[y][x] = getRuneFromColorMap(i)
				}
			}
		}
	}

	sf := Surface{
		Blocks: m,
		Width:  maxC,
		Height: maxR,
		Alpha:  alpha,
	}
	return sf
}

func NewSurfacesFromPngSheet(fn string, r image.Rectangle, alpha bool) []Surface {
	var surfs []Surface

	f, err := os.Open(fn)
	if err != nil {
		//
	}

	img, err := png.Decode(f)
	if err != nil {
		//
	}
	bnd := img.Bounds()
	w := r.Max.X
	i := img.(*image.Paletted)

	for cnt := 0; cnt < bnd.Max.X/w; cnt++ {
		rect := image.Rect(cnt*w, 0, cnt*w+w, r.Max.Y)
		surfs = append(surfs, NewSurfaceFromImage(i.SubImage(rect), alpha))
	}
	return surfs
}

// Clear removes all blocks from a Surface
func (s *Surface) Clear() {
	blocks := make([][]rune, s.Height, s.Height)
	for cnt := 0; cnt < s.Height; cnt++ {
		blocks[cnt] = make([]rune, s.Width, s.Width)
	}
	s.Blocks = blocks
}

// Fill fills an entire surface with a rune
func (s *Surface) Fill(ch rune) {
	for rcnt, r := range s.Blocks {
		for ccnt, _ := range r {
			s.Blocks[rcnt][ccnt] = ch
		}
	}
}

// ConvertToCostume converts a Surface into a Costume usable in a Sprite
func (s Surface) ConvertToCostume() Costume {
	blocks := []*Block{}

	for rcnt := 0; rcnt < len(s.Blocks); rcnt += 2 {
		// XXX - needs to be max(len(m[rcnt]), len(m[rcnt+1]))
		// for ccnt := 0; ccnt < max(len(m[rcnt]), len(m[rcnt+1])); ccnt+=2 {
		for ccnt := 0; ccnt < len(s.Blocks[rcnt]); ccnt += 2 {
			c := 0
			if s.Blocks[rcnt][ccnt] != 0 {
				c += 1
			}
			if len(s.Blocks[rcnt]) > ccnt+1 && s.Blocks[rcnt][ccnt+1] != 0 {
				c += 2
			}
			if len(s.Blocks) > rcnt+1 && s.Blocks[rcnt+1][ccnt] != 0 {
				c += 4
			}
			if len(s.Blocks) > rcnt+1 && len(s.Blocks[rcnt]) > ccnt+1 && s.Blocks[rcnt+1][ccnt+1] == 'X' {
				c += 8
			}

			if (s.Alpha && c > 0) || (!s.Alpha) {
				b := &Block{
					Char: Blocks[c],
					X:    ccnt / 2,
					Y:    rcnt / 2,
				}
				blocks = append(blocks, b)
			}
		}
	}
	return Costume{Blocks: blocks, Width: s.Width / 2}
}

// ConvertToColorCostume converts a Surface into a color Costume usable in a Sprite
func (s Surface) ConvertToColorCostume(bg tm.Attribute) Costume {
	blocks := []*Block{}

	for rcnt := 0; rcnt < len(s.Blocks); rcnt += 2 {
		for ccnt := 0; ccnt < len(s.Blocks[rcnt]); ccnt += 2 {
			var fg tm.Attribute
			obg := bg

			runes := []rune{
				s.Blocks[rcnt][ccnt],
				s.Blocks[rcnt][ccnt+1],
				s.Blocks[rcnt+1][ccnt],
				s.Blocks[rcnt+1][ccnt+1],
			}

			for _, b := range runes {
				if b > 0 && fg == 0 {
					fg = ColorMap[b]
				} else if b != 0 && ColorMap[b] != fg {
					obg = ColorMap[b]
				}
			}

			// if we didn't set a foreground, just skip the block
			if fg == 0 {
				continue
			}

			c := 0
			for cnt, b := range runes {
				if ColorMap[b] == fg {
					c += int(uint(1) << uint(cnt))
				}
			}

			blk := &Block{
				Char: Blocks[c],
				X:    ccnt / 2,
				Y:    rcnt / 2,
				Fg:   tm.Attribute(fg),
				Bg:   tm.Attribute(obg),
			}
			blocks = append(blocks, blk)
		}
	}

	costume := Costume{Blocks: blocks}

	return costume
}

// Blit a Surface onto a Surface
func (s Surface) Blit(t Surface, x, y int) error {
	for rcnt, r := range t.Blocks {
		for ccnt, c := range r {
			if c > 0 {
				s.Point(ccnt+x, rcnt+y, c)
			}
		}
	}
	return nil
}

// Draw a line between two points on a Surface
func (s Surface) Line(x0, y0, x1, y1 int, ch rune) error {
	points := findPointsInLine(x0, y0, x1, y1)
	for _, p := range points {
		s.Point(p.X, p.Y, ch)
	}
	return nil
}

func findPointsInLine(x0, y0, x1, y1 int) []Point {
	var points []Point
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	err := dx + dy
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	for {
		points = append(points, Point{x0, y0})
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
	return points
}

// Draw a rectangle on a Surface
func (s Surface) Rectangle(x0, y0, x1, y1 int, ch rune, fill bool) error {
	if fill {
		xMin := min(x0, x1)
		xMax := max(x0, x1)
		for x := 0; x < xMax-xMin; x++ {
			s.Line(xMin+x, y0, xMin+x, y1, ch)
		}
	} else {
		s.Line(x0, y0, x1, y0, ch)
		s.Line(x1, y0, x1, y1, ch)
		s.Line(x0, y0, x0, y1, ch)
		s.Line(x0, y1, x1, y1, ch)
	}
	return nil
}

// Draw a point on a Surface
func (s Surface) Point(x, y int, ch rune) {
	if y >= 0 && y < len(s.Blocks) {
		if x >= 0 && x < len(s.Blocks[y]) {
			s.Blocks[y][x] = ch
		}
	}
}

// Draw a triangle on a Surface
func (s Surface) Triangle(x0, y0, x1, y1, x2, y2 int, ch rune, fill bool) error {
	if fill {
		points := []Point{
			Point{x0, y0},
			Point{x1, y1},
			Point{x2, y2},
		}

		sort.Slice(points, func(i, j int) bool {
			return points[i].Y < points[j].Y
		})

		pMin := points[0]
		pMid := points[1]
		pMax := points[2]
		pl := findPointsInLine(pMin.X, pMin.Y, pMax.X, pMax.Y)

		var opl []Point

		// don't put in horizontal lines
		if pMin.Y == pMid.Y {
			opl = append(opl, pMid)
		} else {
			opl = append(opl, findPointsInLine(pMin.X, pMin.Y, pMid.X, pMid.Y)...)
		}
		if pMax.Y == pMid.Y {
			opl = append(opl, pMid)
		} else {
			opl = append(opl, findPointsInLine(pMax.X, pMax.Y, pMid.X, pMid.Y)...)
		}

		for _, p := range pl {
			for _, op := range opl {
				if p.Y == op.Y {
					s.Line(p.X, p.Y, op.X, op.Y, ch)
				}
			}
		}
	} else {
		s.Line(x0, y0, x1, y1, ch)
		s.Line(x1, y1, x2, y2, ch)
		s.Line(x2, y2, x0, y0, ch)
	}
	return nil
}

// TexturedTriangle draws a filled triangle on the Surface, sampling colors
// from a texture Surface using barycentric interpolation of the given UV
// coordinates. UV values are normalized (0-1) into tex's pixel space.
// Pixels sampled with a zero rune in the texture are treated as transparent.
// shade scales the brightness of sampled pixels (0.0 is black, 1.0 is
// unshaded); shading is quantized to the closest colour in the xterm palette.
func (s Surface) TexturedTriangle(x0, y0, x1, y1, x2, y2 int, tex Surface, u0, v0, u1, v1, u2, v2, shade float64) {
	if shade <= 0 {
		return
	}

	minX := min(x0, min(x1, x2))
	maxX := max(x0, max(x1, x2))
	minY := min(y0, min(y1, y2))
	maxY := max(y0, max(y1, y2))

	d := (y1-y2)*(x0-x2) + (x2-x1)*(y0-y2)
	if d == 0 {
		return
	}

	for y := minY; y <= maxY; y++ {
		if y < 0 || y >= s.Height {
			continue
		}
		for x := minX; x <= maxX; x++ {
			if x < 0 || x >= s.Width {
				continue
			}
			w0 := float64((y1-y2)*(x-x2)+(x2-x1)*(y-y2)) / float64(d)
			w1 := float64((y2-y0)*(x-x2)+(x0-x2)*(y-y2)) / float64(d)
			w2 := 1.0 - w0 - w1
			if w0 < 0 || w1 < 0 || w2 < 0 {
				continue
			}
			u := w0*u0 + w1*u1 + w2*u2
			v := w0*v0 + w1*v1 + w2*v2
			tx := int(u * float64(tex.Width))
			ty := int(v * float64(tex.Height))
			ch := tex.Pixel(tx, ty)
			if ch != 0 {
				s.Point(x, y, ShadeRune(ch, shade))
			}
		}
	}
}

// ShadeRune applies a brightness factor to a ColorMap rune, returning
// the nearest-colour rune at the scaled brightness. Returns 0 for runes
// not in the ColorMap.
func ShadeRune(ch rune, shade float64) rune {
	attr, ok := ColorMap[ch]
	if !ok {
		return 0
	}
	// attr-1 is the palette index
	pidx := int(attr) - 1
	if pidx < 0 || pidx >= len(palette.Xterm) {
		return ch
	}
	r, g, b, _ := palette.Xterm[pidx].RGBA()
	if shade > 1.0 {
		shade = 1.0
	}
	if shade <= 0 {
		return ch
	}
	c := imageColorRGBA(uint8(r>>8), uint8(g>>8), uint8(b>>8), shade)
	return getRuneFromColorMap(palette.NearestIndex(c))
}

func imageColorRGBA(r, g, b uint8, shade float64) color.RGBA {
	return color.RGBA{
		uint8(float64(r) * shade),
		uint8(float64(g) * shade),
		uint8(float64(b) * shade),
		0xff,
	}
}

// Pixel returns the rune at x, y on the Surface, clamped to the Surface's
// bounds. Returns 0 for empty cells.
func (s Surface) Pixel(x, y int) rune {
	if x < 0 {
		x = 0
	}
	if x >= s.Width {
		x = s.Width - 1
	}
	if y < 0 {
		y = 0
	}
	if y >= s.Height {
		y = s.Height - 1
	}
	if len(s.Blocks) == 0 || len(s.Blocks[y]) == 0 {
		return 0
	}
	return s.Blocks[y][x]
}

// Draw a circle on a Surface
func (s Surface) Circle(xc, yc, r int, ch rune, fill bool) error {
	x := 0
	y := r
	d := 3 - 2*r

	s.drawCircle(xc, yc, x, y, ch, fill)
	for y >= x {
		x++
		if d > 0 {
			y--
			d = d + 4*(x-y) + 10
		} else {
			d = d + 4*x + 6
		}
		s.drawCircle(xc, yc, x, y, ch, fill)
	}

	return nil
}

func (s Surface) drawCircle(xc, yc, x, y int, ch rune, fill bool) {
	if fill {
		s.Line(xc+x, yc+y, xc+x, yc-y, ch)
		s.Line(xc-x, yc+y, xc-x, yc-y, ch)
		s.Line(xc+y, yc+x, xc-y, yc+x, ch)
		s.Line(xc+y, yc-x, xc-y, yc-x, ch)
	} else {
		s.Point(xc+x, yc+y, ch)
		s.Point(xc-x, yc+y, ch)
		s.Point(xc+x, yc-y, ch)
		s.Point(xc-x, yc-y, ch)
		s.Point(xc+y, yc+x, ch)
		s.Point(xc-y, yc+x, ch)
		s.Point(xc+y, yc-x, ch)
		s.Point(xc-y, yc-x, ch)
	}
}

func getRuneFromColorMap(idx int) rune {
	for k, v := range ColorMap {
		if v == tm.Attribute(idx) {
			return k
		}
	}
	ColorMap[rune(idx)] = tm.Attribute(idx)
	return rune(idx)
}
