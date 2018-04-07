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

	n1 := NewWhale(3, 3, sprite.Costume{whale_c0})
	n1.Alpha = 'x'
	n1.VX = 1
	n1.VY = 1

	n2 := NewWhale(70, 3, sprite.Costume{whale_c0})
	n2.Alpha = 'x'
	n2.VX = -1
	n2.VY = 1

	allSprites.Sprites = append(allSprites.Sprites, n1)
	allSprites.Sprites = append(allSprites.Sprites, n2)

mainloop:
	for {
		tm.Clear(tm.ColorDefault, tm.ColorDefault)

		select {
		case ev := <-event_queue:
			if ev.Type == tm.EventKey && ev.Key == tm.KeyEsc {
				break mainloop
			} else if ev.Type == tm.EventResize {
				Width = ev.Width
				Height = ev.Height
			}
		default:
			allSprites.Render()
			allSprites.Update()
			time.Sleep(100 * time.Millisecond)
		}
	}

}
