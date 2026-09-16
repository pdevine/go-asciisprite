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

func main() {
	// XXX - Wait a bit until the terminal is properly initialized
	time.Sleep(500 * time.Millisecond)

	err := tm.Init()
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	Width, Height = tm.Size()

	event_queue := make(chan tm.Event)
	go func() {
		for {
			event_queue <- tm.PollEvent()
		}
	}()

	n1 := NewWhale()
	n2 := NewWhale()

	txt := "Press 'a' to add whales, 'z' to remove them.  'ESC' to quit."
	c := sprite.NewCostume(txt, '~')
	text := sprite.NewBaseSprite(Width/2-len(txt)/2, Height-2, c)

	allSprites.Sprites = append(allSprites.Sprites, text)
	allSprites.Sprites = append(allSprites.Sprites, n1)
	allSprites.Sprites = append(allSprites.Sprites, n2)

mainloop:
	for {
		tm.Clear(tm.ColorDefault, tm.ColorDefault)

		select {
		case ev := <-event_queue:
			if ev.Type == tm.EventKey {
				if ev.Key == tm.KeyEsc {
					break mainloop
				} else if ev.Ch == 'a' {
					w := NewWhale()
					allSprites.Sprites = append(allSprites.Sprites, w)
				} else if ev.Ch == 'z' {
					if len(allSprites.Sprites) > 1 {
						allSprites.Sprites = allSprites.Sprites[:len(allSprites.Sprites)-1]
					}
				}
			} else if ev.Type == tm.EventResize {
				Width = ev.Width
				Height = ev.Height
			}
		default:
			allSprites.Update()
			allSprites.Render()
			time.Sleep(50 * time.Millisecond)
		}
	}

}

const whale_c0 = `xxxxxxxxxxxxxxx##xxxxxxxx.xxxxx
xxxxxxxxx##x##x##xxxxxxx==xxxxx
xxxxxx##x##x##x##xxxxxx===xxxxx
xx/""""""""""""""""\___/x===xxx
x{                      /xx===x
xx\______ o          __/xxxxxxx
xxxx\    \        __/xxxxxxxxxx
xxxxx\____\______/xxxxxxxxxxxxx`

type Whale struct {
	sprite.BaseSprite
	VX int
	VY int
}

func randPos() (int, int) {
	offset := 20
	var x, y int
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	x = r.Intn(Width-2*offset) + offset
	y = r.Intn(Height-2*offset) + offset
	return x, y
}

func randVec() (int, int) {
	var x, y int
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
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

func NewWhale() *Whale {
	s := &Whale{BaseSprite: sprite.BaseSprite{
		Alpha:          'x',
		Height:         0,
		Width:          0,
		Visible:        true,
		Costumes:       []*sprite.Costume{},
		CurrentCostume: 0,
	},
	}
	s.X, s.Y = randPos()
	s.VX, s.VY = randVec()
	s.AddCostume(sprite.NewCostume(whale_c0, 'x'))
	return s
}

func (s *Whale) Update() {
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
