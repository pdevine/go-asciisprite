package main

// PNG sprite demo: dog.png loaded as a two-frame animated costume,
// bouncing around.  'a' adds a dog, 'z' removes one, ESC/q quits.
// Run from the test/ directory: cd test && go run test-png.go

import (
	"math/rand"
	"time"

	sprite "github.com/pdevine/go-asciisprite"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

var allSprites sprite.SpriteGroup
var Width, Height int

type Dog struct {
	sprite.BaseSprite
	VX, VY  int
	Timer   int
	TimeOut int
}

func NewDog() *Dog {
	s := &Dog{
		BaseSprite: sprite.BaseSprite{Visible: true},
		VX:         2,
		VY:         1,
		TimeOut:    10,
	}
	s1 := sprite.NewSurfaceFromPng("testdata/dog.png", true)
	s.BlockCostumes = append(s.BlockCostumes, &s1)
	s2 := sprite.NewSurfaceFromPng("testdata/dog2.png", true)
	s.BlockCostumes = append(s.BlockCostumes, &s2)
	s.Init()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s.X = r.Intn(Width)
	s.Y = r.Intn(Height)
	if r.Intn(2) == 0 {
		s.VX = -s.VX
	}
	if r.Intn(2) == 0 {
		s.VY = -s.VY
	}
	return s
}

func (s *Dog) Update() {
	s.Timer++
	if s.Timer >= s.TimeOut {
		s.Timer = 0
		s.CurrentCostume++
		if s.CurrentCostume >= len(s.BlockCostumes) {
			s.CurrentCostume = 0
		}
	}
	s.X += s.VX
	s.Y += s.VY
	if s.X <= 0 || s.X >= Width-s.Width {
		s.VX = -s.VX
	}
	if s.Y <= 0 || s.Y >= Height-s.Height {
		s.VY = -s.VY
	}
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

	event_queue := make(chan tm.Event)
	go func() {
		for {
			event_queue <- tm.PollEvent()
		}
	}()

	allSprites.Init(Width, Height, true)
	allSprites.BlockMode = true
	allSprites.Background = tm.Attribute(178)
	allSprites.Sprites = append(allSprites.Sprites, NewDog())

mainloop:
	for {
		tm.Clear(tm.Attribute(178), tm.Attribute(178))

		select {
		case ev := <-event_queue:
			if ev.Type == tm.EventKey || ev.Type == tm.EventKeyPress {
				if ev.Key == tm.KeyEsc || ev.Ch == 'q' {
					break mainloop
				} else if ev.Ch == 'a' {
					allSprites.Sprites = append(allSprites.Sprites, NewDog())
				} else if ev.Ch == 'z' {
					if len(allSprites.Sprites) > 1 {
						allSprites.Sprites = allSprites.Sprites[:len(allSprites.Sprites)-1]
					}
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
