package main

import (
	"time"
	"math/rand"

	sprite "github.com/pdevine/go-asciisprite"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

var allSprites sprite.SpriteGroup
var Width int
var Height int

const invader_c0 = `  X     X
 XXXXXXXXX
 XX XXX XX
 XXXXXXXXX
 X X   X X
X X   X X`

const invader_c1 = `  X     X
 XXXXXXXXX
 XX XXX XX
 XXXXXXXXX
 X X   X X
  X X   X X`


type Invader struct {
	sprite.BaseSprite
	VX      int
	VY      int
	Timer   int
	TimeOut int
}
func randPos() (int, int) {
	offset := 20
	var x, y int
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
        x = r.Intn(Width-2*offset) + offset
        y = r.Intn(Height-2*offset) + offset
        return x, y
}

func randVec() (int, int) {
        var x, y int
        r := rand.New(rand.NewSource(time.Now().UnixNano()))
        n := r.Intn(2)
        x = n
        if x == 0 {
                x = -1
        }

        n = r.Intn(2)
        y = n
        if y == 0 {
                y = -1
        }
        return x, y
}


func NewInvader() *Invader {
	s := &Invader{BaseSprite: sprite.BaseSprite{
		Visible: true},
		TimeOut: 10,
	}

	s1 := sprite.NewSurfaceFromPng("dog.png")
	s.BlockCostumes = append(s.BlockCostumes, &s1)
	s2 := sprite.NewSurfaceFromPng("dog2.png")
	s.BlockCostumes = append(s.BlockCostumes, &s2)

	s.X, s.Y = randPos()
	s.VX, s.VY = randVec()

	return s
}

func (s *Invader) Update() {
	s.Timer += 1
	if s.Timer >= s.TimeOut {
		s.Timer = 0
		s.CurrentCostume += 1
		if s.CurrentCostume >= len(s.BlockCostumes) {
			s.CurrentCostume = 0
		}
	}

	s.X = s.X + s.VX
	s.Y = s.Y + s.VY

	if s.X < 0 {
		s.VX = 1
	}
	if s.X > Width-s.Width {
		s.VX = -1
	}
	if s.Y >= Height-s.Height {
		s.VY = -1
	}
	if s.Y <= 0 {
		s.VY = 1
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
	Width = w*2
	Height = h*2

	event_queue := make(chan tm.Event)
	go func() {
		for {
			event_queue <- tm.PollEvent()
		}
	}()

	i := NewInvader()

	allSprites.Init(Width, Height, true)
	allSprites.BlockMode = true
	allSprites.Sprites = append(allSprites.Sprites, i)

mainloop:
	for {
		tm.Clear(tm.ColorDefault, tm.ColorDefault)

		select {
		case ev := <-event_queue:
			if ev.Type == tm.EventKey {
				if ev.Key == tm.KeyEsc {
					break mainloop
				} else if ev.Ch == 'a' {
					i := NewInvader()
					allSprites.Sprites = append(allSprites.Sprites, i)
				} else if ev.Ch == 'z' {
					if len(allSprites.Sprites) > 1 {
						allSprites.Sprites = allSprites.Sprites[:len(allSprites.Sprites)-1]
					}
				}
			} else if ev.Type == tm.EventResize {
				Width = ev.Width*2
				Height = ev.Height*2
				allSprites.Resize(Width, Height)
			}
		default:
			allSprites.Update()
			allSprites.Render()
			time.Sleep(50 * time.Millisecond)
		}
	}

}
