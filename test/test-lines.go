package main

import (
	"math/rand"
	"time"

	sprite "github.com/pdevine/go-asciisprite"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

var allSprites sprite.SpriteGroup
var Width int
var Height int
var Rand *rand.Rand

type Lines struct {
	sprite.BaseSprite
	timer   int
	timeOut int
}

func NewLines() *Lines {
	s := &Lines{BaseSprite: sprite.BaseSprite{
		Visible: true},
		timeOut: 100,
	}

	surf := sprite.NewSurface(Width, Height, false)
	surf.Line(2, 2, Width-2, Height-2, 'X')

	s.BlockCostumes = append(s.BlockCostumes, &surf)

	return s
}

func (s *Lines) Update() {
	s.timer++
	if s.timer >= s.timeOut {
		s.timer = 0
		surf := sprite.NewSurface(Width, Height, false)
		s.BlockCostumes[0] = &surf
	}
	x1 := Rand.Intn(Width)
	y1 := Rand.Intn(Height)
	x2 := Rand.Intn(Width)
	y2 := Rand.Intn(Height)
	s.BlockCostumes[0].Line(x1, y1, x2, y2, 'X')
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

	i := NewLines()

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
