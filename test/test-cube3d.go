package main

import (
	"math"
	"math/rand"
	"sort"
	"time"

	sprite "github.com/pdevine/go-asciisprite"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

var allSprites sprite.SpriteGroup
var Width int
var Height int
var Rand *rand.Rand

type Triangle struct {
	A     *Point3D
	B     *Point3D
	C     *Point3D
	Color rune
	surf  *sprite.Surface
}

func NewTriangle(a, b, c *Point3D, ch rune) *Triangle {
	t := &Triangle{
		A:     a,
		B:     b,
		C:     c,
		Color: ch,
	}
	return t
}

func (t *Triangle) isBackFace() bool {
	cax := t.C.ScreenX() - t.A.ScreenX()
	cay := t.C.ScreenY() - t.A.ScreenY()

	bcx := t.B.ScreenX() - t.C.ScreenX()
	bcy := t.B.ScreenY() - t.C.ScreenY()

	return cax * bcy > cay * bcx
}

func (t *Triangle) Depth() float64 {
	zPos := math.Min(t.A.Z, t.B.Z)
	zPos = math.Min(zPos, t.C.Z)
	return zPos
}

func (t *Triangle) Draw(surf *sprite.Surface) {
	if t.isBackFace() {
		return
	}

	surf.Triangle(
		t.A.ScreenX(), t.A.ScreenY(),
		t.B.ScreenX(), t.B.ScreenY(),
		t.C.ScreenX(), t.C.ScreenY(),
		t.Color, true)
}

type Point3D struct {
	X     float64
	Y     float64
	Z     float64
	cX    float64
	cY    float64
	cZ    float64
	fl    float64
	vpX   float64
	vpY   float64
	scale float64
}

func NewPoint3D(x, y, z float64) *Point3D {
	p := &Point3D{
		fl:    250.0,
		scale: 1.0,
		X:     x,
		Y:     y,
		Z:     z,
	}
	return p
}

func (p *Point3D) SetVanishingPoint(vpX, vpY int) {
	p.vpX = float64(vpX)
	p.vpY = float64(vpY)
}

func (p *Point3D) SetCenter(cX, cY, cZ float64) {
	p.cX = cX
	p.cY = cY
	p.cZ = cZ
}

func (p *Point3D) ScreenX() int {
	p.scale = p.fl / (p.fl + p.Z + p.cZ)
	return int(math.Round(p.vpX + (p.cX+p.X)*p.scale))
}

func (p *Point3D) ScreenY() int {
	p.scale = p.fl / (p.fl + p.Z + p.cZ)
	return int(math.Round(p.vpY + (p.cY+p.Y)*p.scale))
}

func (p *Point3D) RotateX(angleX float64) {
	cosX := math.Cos(angleX)
	sinX := math.Sin(angleX)

	y := (p.Y * cosX) - (p.Z * sinX)
	z := (p.Z * cosX) + (p.Y * sinX)

	p.Y = y
	p.Z = z
}

func (p *Point3D) RotateY(angleY float64) {
	cosY := math.Cos(angleY)
	sinY := math.Sin(angleY)

	x := (p.X * cosY) - (p.Z * sinY)
	z := (p.Z * cosY) + (p.X * sinY)

	p.X = x
	p.Z = z
}

func (p *Point3D) RotateZ(angleZ float64) {
	cosZ := math.Cos(angleZ)
	sinZ := math.Sin(angleZ)

	x := (p.X * cosZ) - (p.Y * sinZ)
	y := (p.Y * cosZ) + (p.X * sinZ)

	p.X = x
	p.Y = y
}

type Square3D struct {
	sprite.BaseSprite
	points    []*Point3D
	triangles []*Triangle
}

func NewSquare3D() *Square3D {
	s := &Square3D{BaseSprite: sprite.BaseSprite{
		Visible: true},
	}

	s.points = []*Point3D{
		NewPoint3D(-30, -30, -30),
		NewPoint3D( 30, -30, -30),
		NewPoint3D( 30,  30, -30),
		NewPoint3D(-30,  30, -30),

		NewPoint3D(-30, -30,  30),
		NewPoint3D( 30, -30,  30),
		NewPoint3D( 30,  30,  30),
		NewPoint3D(-30,  30,  30),
	}

	for _, p := range s.points {
		p.SetVanishingPoint(Width/2, Height/2)
		p.SetCenter(0, 0,  10)
	}

	s.triangles = []*Triangle{
		// front
		NewTriangle(s.points[0], s.points[1], s.points[2], 'X'),
		NewTriangle(s.points[0], s.points[2], s.points[3], 'X'),

		// top
		NewTriangle(s.points[0], s.points[5], s.points[1], 'b'),
		NewTriangle(s.points[0], s.points[4], s.points[5], 'b'),

		// back
		NewTriangle(s.points[4], s.points[6], s.points[5], 'N'),
		NewTriangle(s.points[4], s.points[7], s.points[6], 'N'),

		// bottom
		NewTriangle(s.points[3], s.points[2], s.points[6], 'G'),
		NewTriangle(s.points[3], s.points[6], s.points[7], 'G'),

		// right
		NewTriangle(s.points[1], s.points[5], s.points[6], 'B'),
		NewTriangle(s.points[1], s.points[6], s.points[2], 'B'),

		// left
		NewTriangle(s.points[4], s.points[0], s.points[3], 'p'),
		NewTriangle(s.points[4], s.points[3], s.points[7], 'p'),

	}

	surf := sprite.NewSurface(Width, Height, false)
	s.BlockCostumes = append(s.BlockCostumes, &surf)

	return s
}

func (s *Square3D) Update() {
	angleX := 0.05
	angleY := 0.1
	angleZ := 0.05

	for _, p := range s.points {
		p.RotateX(angleX)
		p.RotateY(angleY)
		p.RotateZ(angleZ)
	}

	sort.Slice(s.triangles, func(i, j int) bool {
		return s.triangles[i].Depth() > s.triangles[j].Depth()
	})

	surf := sprite.NewSurface(Width, Height, false)
	for _, t := range s.triangles {
		t.Draw(&surf)
	}

	s.BlockCostumes[0] = &surf
}

func setPalette() {
	sprite.ColorMap['b'] = tm.Attribute(99)
	sprite.ColorMap['N'] = tm.ColorNavy
	sprite.ColorMap['G'] = tm.ColorGray
	sprite.ColorMap['B'] = tm.Attribute(111)
	sprite.ColorMap['p'] = tm.Attribute(131)
}

func main() {
	// XXX - Wait a bit until the terminal is properly initialized
	time.Sleep(500 * time.Millisecond)

	err := tm.Init()
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	w, h := tm.Size()
	Width = w * 2
	Height = h * 2
	Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	setPalette()

	event_queue := make(chan tm.Event)
	go func() {
		for {
			event_queue <- tm.PollEvent()
		}
	}()

	i := NewSquare3D()

	allSprites.Init(Width, Height, true)
	allSprites.BlockMode = true
	allSprites.Background = tm.Attribute(178)
	allSprites.Sprites = append(allSprites.Sprites, i)

mainloop:
	for {
		tm.Clear(tm.Attribute(178), tm.Attribute(178))

		select {
		case ev := <-event_queue:
			if ev.Type == tm.EventKey {
				if ev.Key == tm.KeyEsc {
					break mainloop
				}
			} else if ev.Type == tm.EventResize {
				Width = ev.Width * 2
				Height = ev.Height * 2
				allSprites.Resize(Width, Height)
			}
		default:
			allSprites.Update()
			allSprites.Render()
			time.Sleep(50 * time.Millisecond)
		}
	}
}
