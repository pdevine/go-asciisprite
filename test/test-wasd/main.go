package main

// WASD movement demo for the kitty keyboard protocol support.
//
// Hold W/A/S/D (or arrows) to move the block; release to stop.
// Press +/- to resize it, q or Esc to quit.
//
// On terminals without kitty keyboard protocol support the demo falls
// back to legacy key events: movement happens one step per keypress and
// there is no way to track held keys.  A banner shows which mode is
// active.

import (
	"time"

	sprite "github.com/pdevine/go-asciisprite"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

var allSprites sprite.SpriteGroup
var Width int
var Height int
var enhanced bool

var player *Block

type Block struct {
	sprite.BaseSprite
	Size int
}

func (b *Block) costume() *sprite.Surface {
	s := sprite.NewSurface(b.Size*2, b.Size, true)
	s.Fill('R') // ColorMap red
	return &s
}

var held = map[interface{}]bool{}

func keyDir() (dx, dy int) {
	if held['w'] || held[tm.KeyArrowUp] {
		dy--
	}
	if held['s'] || held[tm.KeyArrowDown] {
		dy++
	}
	if held['a'] || held[tm.KeyArrowLeft] {
		dx--
	}
	if held['d'] || held[tm.KeyArrowRight] {
		dx++
	}
	return dx, dy
}

func (b *Block) Update() {
	dx, dy := keyDir()
	b.X += dx * 2
	b.Y += dy
	if b.X < 0 {
		b.X = 0
	}
	if b.Y < 2 {
		b.Y = 2
	}
	if b.X > Width-b.Size*2 {
		b.X = Width - b.Size*2
	}
	if b.Y > Height-b.Size {
		b.Y = Height - b.Size
	}
}

func step(ch rune, k tm.Key) {
	switch {
	case ch == 'w' || k == tm.KeyArrowUp:
		player.Y--
	case ch == 's' || k == tm.KeyArrowDown:
		player.Y++
	case ch == 'a' || k == tm.KeyArrowLeft:
		player.X -= 2
	case ch == 'd' || k == tm.KeyArrowRight:
		player.X += 2
	}
}

func resize(delta int) {
	player.Size += delta
	if player.Size < 2 {
		player.Size = 2
	}
	if player.Size > 20 {
		player.Size = 20
	}
	player.BlockCostumes[0] = player.costume()
	player.Init()
}

func main() {
	// XXX - Wait a bit until the terminal is properly initialized
	time.Sleep(500 * time.Millisecond)

	flags, err := tm.InitEnhancedKeys()
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	enhanced = flags != 0

	w, h := tm.Size()
	Width = w * 2
	Height = h * 2

	event_queue := make(chan tm.Event)
	go func() {
		for {
			event_queue <- tm.PollEvent()
		}
	}()

	status := "legacy keys: one step per press"
	if enhanced {
		status = "enhanced keys (kitty protocol): hold WASD/arrows to move"
	}

	player = &Block{BaseSprite: sprite.BaseSprite{Visible: true}, Size: 6}
	player.BlockCostumes = append(player.BlockCostumes, player.costume())
	player.Init()
	player.X, player.Y = Width/2, Height/2

	txt := sprite.NewCostume(status+"  |  +/- resize, q/ESC quit", 'x')
	text := sprite.NewBaseSprite(2, 0, txt)

	allSprites.Init(Width, Height, true)
	allSprites.BlockMode = true
	allSprites.Background = tm.Attribute(16)
	allSprites.Sprites = append(allSprites.Sprites, player, text)

mainloop:
	for {
		tm.Clear(tm.Attribute(16), tm.Attribute(16))

		select {
		case ev := <-event_queue:
			switch ev.Type {
			case tm.EventKeyPress:
				if ev.Key == tm.KeyEsc || ev.Ch == 'q' {
					break mainloop
				}
				switch {
				case ev.Ch == '+' || ev.Ch == '=':
					resize(1)
				case ev.Ch == '-' || ev.Ch == '_':
					resize(-1)
				default:
					if ev.Ch != 0 {
						held[ev.Ch] = true
					}
					if ev.Key != 0 {
						held[ev.Key] = true
					}
				}
			case tm.EventKeyRepeat:
				// already marked held; nothing to do
			case tm.EventKeyRelease:
				if ev.Ch != 0 {
					delete(held, ev.Ch)
				}
				if ev.Key != 0 {
					delete(held, ev.Key)
				}
			case tm.EventKey:
				// legacy fallback: press = step
				if ev.Key == tm.KeyEsc || ev.Ch == 'q' {
					break mainloop
				}
				switch {
				case ev.Ch == '+' || ev.Ch == '=':
					resize(1)
				case ev.Ch == '-' || ev.Ch == '_':
					resize(-1)
				default:
					step(ev.Ch, ev.Key)
				}
			case tm.EventResize:
				Width = ev.Width * 2
				Height = ev.Height * 2
				allSprites.Resize(Width, Height)
			}
		default:
			allSprites.Update()
			allSprites.Render()
			time.Sleep(20 * time.Millisecond)
		}
	}
}
