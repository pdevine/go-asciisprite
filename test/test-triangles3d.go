package main

import (
	"math"
	"math/rand"
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

func (t *Triangle) Draw(surf *sprite.Surface) {
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
		NewPoint3D(-25, -35, 10),
		NewPoint3D( 15, -35, 10),
		NewPoint3D(-25, -25, 10),
		NewPoint3D(-15, -25, 10),
		NewPoint3D( 15, -25, 10),

		// 5

		NewPoint3D( 25, -25, 10),
		NewPoint3D(-35, -25, 10),
		NewPoint3D( 35, -25, 10),
		NewPoint3D(-35, -15, 10),
		NewPoint3D( 35, -15, 10),

		// 10

		NewPoint3D(-45, -15, 10),
		NewPoint3D(-25, -15, 10),
		NewPoint3D(-15, -15, 10),
		NewPoint3D( 15, -15, 10),
		NewPoint3D( 25, -15, 10),

		// 15

		NewPoint3D(-45,  -5, 10),
		NewPoint3D(-35,  -5, 10),
		NewPoint3D(-25,  -5, 10),
		NewPoint3D(-15,  -5, 10),
		NewPoint3D( 15,  -5, 10),

		// 20

		NewPoint3D( 25,  -5, 10),
		NewPoint3D( 35,  -5, 10),
		NewPoint3D( 45,  -5, 10),
		NewPoint3D(-55,  -5, 10),
		NewPoint3D( 55,  -5, 10),

		// 25

		NewPoint3D(-55,   5, 10),
		NewPoint3D(-35,   5, 10),
		NewPoint3D( 35,   5, 10),
		NewPoint3D( 55,   5, 10),
		NewPoint3D(-35,  15, 10),

		// 30

		NewPoint3D(-25,  15, 10),
		NewPoint3D( 25,  15, 10),
		NewPoint3D( 35,  15, 10),
		NewPoint3D(-55,  25, 10),
		NewPoint3D(-45,  25, 10),

		// 35

		NewPoint3D(-35,  25, 10),
		NewPoint3D(-25,  25, 10),
		NewPoint3D( -5,  25, 10),
		NewPoint3D(  5,  25, 10),
		NewPoint3D( 25,  25, 10),

		// 40

		NewPoint3D( 35,  25, 10),
		NewPoint3D( 45,  25, 10),
		NewPoint3D( 55,  25, 10),
		NewPoint3D(-25,  35, 10),
		NewPoint3D( -5,  35, 10),

		// 45

		NewPoint3D(  5,  35, 10),
		NewPoint3D( 25,  35, 10),
		NewPoint3D(-45,   5, 10),
		NewPoint3D( 45,   5, 10),
		NewPoint3D(-15, -35, 10),

		// 50

		NewPoint3D( 25, -35, 10),
		NewPoint3D( 45, -15, 10),
	}

	for _, p := range s.points {
		p.SetVanishingPoint(Width/2, Height/2)
		p.SetCenter(0, 0, 100)
	}

	s.triangles = []*Triangle{
		NewTriangle(s.points[0], s.points[49], s.points[3], 'X'),
		NewTriangle(s.points[0], s.points[3], s.points[2], 'X'),
		NewTriangle(s.points[1], s.points[50], s.points[5], 'X'),
		NewTriangle(s.points[1], s.points[5], s.points[4], 'X'),

		NewTriangle(s.points[6], s.points[7], s.points[9], 'X'),
		NewTriangle(s.points[6], s.points[9], s.points[8], 'X'),
		NewTriangle(s.points[10], s.points[11], s.points[17], 'X'),
		NewTriangle(s.points[10], s.points[17], s.points[15], 'X'),
		NewTriangle(s.points[12], s.points[13], s.points[19], 'X'),
		NewTriangle(s.points[12], s.points[19], s.points[18], 'X'),
		NewTriangle(s.points[14], s.points[51], s.points[22], 'X'),
		NewTriangle(s.points[14], s.points[22], s.points[20], 'X'),

		NewTriangle(s.points[23], s.points[16], s.points[26], 'X'),
		NewTriangle(s.points[23], s.points[26], s.points[25], 'X'),
		NewTriangle(s.points[16], s.points[21], s.points[32], 'X'),
		NewTriangle(s.points[16], s.points[32], s.points[29], 'X'),
		NewTriangle(s.points[21], s.points[24], s.points[28], 'X'),
		NewTriangle(s.points[21], s.points[28], s.points[27], 'X'),

		NewTriangle(s.points[25], s.points[47], s.points[34], 'X'),
		NewTriangle(s.points[25], s.points[34], s.points[33], 'X'),
		NewTriangle(s.points[29], s.points[30], s.points[36], 'X'),
		NewTriangle(s.points[29], s.points[36], s.points[35], 'X'),
		NewTriangle(s.points[31], s.points[32], s.points[40], 'X'),
		NewTriangle(s.points[31], s.points[40], s.points[39], 'X'),
		NewTriangle(s.points[48], s.points[28], s.points[42], 'X'),
		NewTriangle(s.points[48], s.points[42], s.points[41], 'X'),

		NewTriangle(s.points[36], s.points[37], s.points[44], 'X'),
		NewTriangle(s.points[36], s.points[44], s.points[43], 'X'),
		NewTriangle(s.points[38], s.points[39], s.points[46], 'X'),
		NewTriangle(s.points[38], s.points[46], s.points[45], 'X'),
	}

	surf := sprite.NewSurface(Width, Height, false)
	s.BlockCostumes = append(s.BlockCostumes, &surf)

	return s
}

func (s *Square3D) Update() {
	angleX := 0.0
	angleY := 0.1
	angleZ := 0.0

	for _, p := range s.points {
		p.RotateX(angleX)
		p.RotateY(angleY)
		p.RotateZ(angleZ)
	}

	surf := sprite.NewSurface(Width, Height, false)
	for _, t := range s.triangles {
		t.Draw(&surf)
	}

	s.BlockCostumes[0] = &surf
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
