```
    .-_'''-.       ,-----.                    {\    /}
   '_( )_   \    .'  .-,  '.                   >`()'<
  |(_ o _)|  '  / ,-.|  \ _ \                 {@\/|\ @}
  . (_,_)/___| ;  \  '_ /  | :  _ _    _ _     `~'|/~'
  |  |  .-----.|  _`,/ \ _/  | ( ' )--( ' )      //
  '  \  '-   .': (  '\_/ \   ;(_{;}_)(_{;}_)     \\
   \  `-'`   |  \ `"/  \  ) /  (_,_)--(_,_)      ''
    \        /   '. \_/``".'                 
     `'-...-'      '-----'                   
  
      ____       .-'''-.     _______  .-./`) .-./`)
    .'  __ `.   / _     \   /   __  \ \ .-.')\ .-.')
   /   '  \  \ (`' )/`--'  | ,_/  \__)/ `-' \/ `-' \
   |___|  /  |(_ o _).   ,-./  )       `-'`"` `-'`"`
      _.-`   | (_,_). '. \  '_ '`)     .---.  .---.
   .'   _    |.---.  \  : > (_)  )  __ |   |  |   |
   |  _( )_  |\    `-'  |(  .  .-'_/  )|   |  |   |
   \ (_ o _) / \       /  `-'`-'     / |   |  |   |
    '.(_,_).'   `-...-'     `._____.'  '---'  '---'

        .-'''-. .-------. .-------.   .-./`) ,---------.    .-''-.   
       / _     \\  _(`)_ \|  _ _   \  \ .-.')\          \ .'_ _   \  
      (`' )/`--'| (_ o._)|| ( ' )  |  / `-' \ `--.  ,---'/ ( ` )   ' 
      _ o _).   |  (_,_) /|(_ o _) /   `-'`"`    |   \  . (_ o _)  | 
      (_,_). '. |   '-.-' | (_,_).' __ .---.     :_ _:  |  (_,_)___| 
     ,---.  \  :|   |     |  |\ \  |  ||   |     (_I_)  '  \   .---. 
     \    `-'  ||   |     |  | \ `'   /|   |    (_(=)_)  \  `-'    / 
      \       / /   )     |  |  \    / |   |     (_I_)    \       /  
       `-...-'  `---'     ''-'   `'-'  '---'     '---'     `'-..-'   

```

A simple golang sprite library for animating ASCII and Unicode art

***What is a sprite?*** A sprite is a two-dimensional object which you can use in videogames and animations.

## Features


 * Easy to use routines for creating and manipulating text sprites

 * Block interpolation for creating unicode "pixel" graphics

 * Simple ascii "fonts" for creating block based text

 * Event system for triggering sprite events


## Usage

```go
package main

import (
        "time"

        sprite "github.com/pdevine/go-asciisprite"
        tm "github.com/pdevine/go-asciisprite/termbox"
)

var allSprites sprite.SpriteGroup
var Width int
var Height int

const cow_c0 = `
^__^xxxxxxxxxxxx
(oo)\_______
(__)\       )\/\
xxxx||----w |xxx
xxxx||xxxxx||xxx
`

type Cow struct {
	sprite.BaseSprite
	VX int
	VY int
}

func NewCow() *Cow {
	cow := &Cow{BaseSprite: sprite.BaseSprite{
		X: 5,
		Y: 5,
		Visible: true},
		VX: 1,
		VY: 1,
	}
	cow.AddCostume(sprite.NewCostume(cow_c0, 'x'))

	return cow
}

func (cow *Cow) Update() {
	cow.X += cow.VX
	cow.Y += cow.VY
	if cow.X < 0 {
		cow.VX = 1
	} else if cow.X > Width - cow.Width {
		cow.VX = -1
	}
	if cow.Y < 0 {
		cow.VY = 1
	} else if cow.Y > Height - cow.Height {
		cow.VY = -1
	}
}

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

	cow := NewCow()
	c := sprite.NewCostume("Press 'ESC' to quit", '~')
	text := sprite.NewBaseSprite(Width/2-c.Width/2, Height-2, c)

	allSprites.Sprites = append(allSprites.Sprites, cow)
	allSprites.Sprites = append(allSprites.Sprites, text)

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
			time.Sleep(60 * time.Millisecond)
		}
	}
}
```

## Featured Projects

 * [Tetromino Elektronica](https://github.com/pdevine/tetromino) A Falling Tetromino Game

 * [Docker Doodle](https://github.com/docker/doodle) ASCII doodles with Moby the Whale



