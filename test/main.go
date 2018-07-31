package main

import (
	"time"

	sprite "github.com/pdevine/go-asciisprite"

	tm "github.com/gdamore/tcell/termbox"
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
