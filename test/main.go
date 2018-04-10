package main

import (
	"time"

	"./sprite"

	tm "github.com/nsf/termbox-go"
)

var allSprites sprite.SpriteGroup
var Width int
var Height int

func main() {
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
					if len(allSprites.Sprites) > 0 {
						allSprites.Sprites = allSprites.Sprites[:len(allSprites.Sprites)-1]
					}
				}
			} else if ev.Type == tm.EventResize {
				Width = ev.Width
				Height = ev.Height
			}
		default:
			allSprites.Render()
			allSprites.Update()
			time.Sleep(50 * time.Millisecond)
		}
	}

}
