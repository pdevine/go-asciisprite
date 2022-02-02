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

	p.Y = (p.Y * cosX) - (p.Z * sinX)
	p.Z = (p.Z * cosX) + (p.Y * sinX)
}

func (p *Point3D) RotateY(angleY float64) {
	cosY := math.Cos(angleY)
	sinY := math.Sin(angleY)

	p.X = (p.X * cosY) - (p.Z * sinY)
	p.Z = (p.Z * cosY) + (p.X * sinY)
}

func (p *Point3D) RotateZ(angleZ float64) {
	cosZ := math.Cos(angleZ)
	sinZ := math.Sin(angleZ)

	p.X = (p.X * cosZ) - (p.Y * sinZ)
	p.Y = (p.Y * cosZ) + (p.X * sinZ)
}

type Square3D struct {
	sprite.BaseSprite
	points []*Point3D
}

func NewSquare3D() *Square3D {
	s := &Square3D{BaseSprite: sprite.BaseSprite{
		Visible: true},
	}

	s.points = []*Point3D{
		NewPoint3D(-10, -10, 50),
		NewPoint3D(10, -10, 50),
		NewPoint3D(10, 10, 50),
		NewPoint3D(-10, 10, 50),
	}

	for _, p := range s.points {
		p.SetVanishingPoint(Width/2, Height/2)
	}

	surf := sprite.NewSurface(Width, Height, false)
	s.BlockCostumes = append(s.BlockCostumes, &surf)

	return s
}

func (s *Square3D) Update() {
	angleX := 0.01
	angleY := 0.05

	for _, p := range s.points {
		p.RotateX(angleX)
		p.RotateY(angleY)
	}

	surf := sprite.NewSurface(Width, Height, false)
	for cnt := 0; cnt < len(s.points); cnt++ {
		if cnt == len(s.points)-1 {
			c := s.points[cnt]
			n := s.points[0]
			surf.Line(c.ScreenX(), c.ScreenY(), n.ScreenX(), n.ScreenY(), 'X')
		} else {
			c := s.points[cnt]
			n := s.points[cnt+1]
			surf.Line(c.ScreenX(), c.ScreenY(), n.ScreenX(), n.ScreenY(), 'X')
		}
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
