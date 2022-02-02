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
				// we only support 256 colour mode, so get the index from the palette
				// and create an entry in our colour map if needed
				i := palette.Index(color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
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
