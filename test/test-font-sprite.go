package main

import (
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

	f := sprite.NewJRFont()
	s := &sprite.BaseSprite{
		Visible: true,
		X:       10,
		Y:       5,
	}
        s.AddCostume(sprite.Convert(f.BuildString("Hello, World!")))

	f2 := sprite.NewJRSMFont()
	s2 := &sprite.BaseSprite{
		Visible: true,
		X:       10,
		Y:       15,
	}
	s2.AddCostume(sprite.Convert(f2.BuildString("the quick brown fox jumps over the lazy dog")))

	f3 := sprite.NewPakuFont()
	s3 := &sprite.BaseSprite{
		Visible: true,
		X:       10,
		Y:       20,
	}
        s3.AddCostume(sprite.Convert(f3.BuildString("The quick brown fox jumps over the lazy dog")))

	s4 := &sprite.BaseSprite{
		Visible: true,
		X:       10,
		Y:       25,
	}
        s4.AddCostume(sprite.Convert(f2.BuildString("THE QUICK BROWN FOX JUMPS OVER THE LAZY DOG")))

	f4 := sprite.New90sFont()
	s5 := &sprite.BaseSprite{
		Visible: true,
		X:       10,
		Y:       30,
	}
        s5.AddCostume(sprite.Convert(f4.BuildString("THE QUICK BROWN FOX JUMPS OVER THE LAZY DOG")))

	s6 := &sprite.BaseSprite{
		Visible: true,
		X:       10,
		Y:       35,
	}
        s6.AddCostume(sprite.Convert(f4.BuildString("the quick brown fox jumps over the lazy dog")))


	txt := "Press 'ESC' to quit."
	c := sprite.NewCostume(txt, '~')
	text := sprite.NewBaseSprite(Width/2-len(txt)/2, Height-2, c)

	allSprites.Sprites = append(allSprites.Sprites, text)
	allSprites.Sprites = append(allSprites.Sprites, s)
	allSprites.Sprites = append(allSprites.Sprites, s2)
	allSprites.Sprites = append(allSprites.Sprites, s3)
	allSprites.Sprites = append(allSprites.Sprites, s4)
	allSprites.Sprites = append(allSprites.Sprites, s5)
	allSprites.Sprites = append(allSprites.Sprites, s6)

mainloop:
	for {
		tm.Clear(tm.ColorDefault, tm.ColorDefault)

		select {
		case ev := <-event_queue:
			if ev.Type == tm.EventKey {
				if ev.Key == tm.KeyEsc {
					break mainloop
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
