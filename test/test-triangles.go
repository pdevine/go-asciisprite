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
	solid   bool
}

func NewLines() *Lines {
	s := &Lines{BaseSprite: sprite.BaseSprite{
		Visible: true},
		timeOut: 100,
	}

	surf := sprite.NewSurface(Width, Height, false)

	s.BlockCostumes = append(s.BlockCostumes, &surf)

	return s
}

func (s *Lines) Update() {
	s.timer++
	if s.timer >= s.timeOut {
		s.timer = 0
		surf := sprite.NewSurface(Width, Height, false)
		s.BlockCostumes[0] = &surf
		s.solid = Rand.Intn(2) == 1
	}
	x0 := Rand.Intn(Width)
	y0 := Rand.Intn(Height)
	x1 := Rand.Intn(Width)
	y1 := Rand.Intn(Height)
	x2 := Rand.Intn(Width)
	y2 := Rand.Intn(Height)
	cSlice := []rune{'y', 'B', 'w', 'b', 'N', 'G'}
	s.BlockCostumes[0].Triangle(x0, y0, x1, y1, x2, y2, cSlice[rand.Intn(len(cSlice))], s.solid)
}

func setPalette() {
	sprite.ColorMap['y'] = tm.ColorYellow
	sprite.ColorMap['B'] = tm.ColorBlack
	sprite.ColorMap['w'] = tm.ColorWhite
	sprite.ColorMap['b'] = tm.ColorBlue
	sprite.ColorMap['N'] = tm.ColorNavy
	sprite.ColorMap['G'] = tm.ColorGray
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
