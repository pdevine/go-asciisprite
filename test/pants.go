package main

import (
	"math/rand"
	"time"

	sprite "github.com/pdevine/go-asciisprite"

	tm "github.com/gdamore/tcell/termbox"
)

var allSprites sprite.SpriteGroup
var Width int
var Height int

const confetti_c0 = `.---,
|   |
'---'`

const confetti_c1 = `  .
 / \
 \ /
  '`
const pants_c0 = `
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx,---.,---.,---.
,------.xxx,---.xx,--.xx,--.,--------.x,---.xx|   ||   ||   |
|  .--. 'x/  O  \x|  ,'.|  |'--.  .--''   .-'x|  .'|  .'|  .'
|  '--' ||  .-.  ||  |' '  |xxx|  |xxx` + "`" + `.  ` + "`" + `-.x|  |x|  |x|  |x
|  | --'x|  |x|  ||  |x` + "`" + `   |xxx|  |xxx.-'    |` + "`" + `--'x` + "`" + `--'x` + "`" + `--'x
` + "`" + `--'xxxxx` + "`" + `--'x` + "`" + `--'` + "`" + `--'xx` + "`" + `--'xxx` + "`" + `--'xxx` + "`" + `-----' .--.x.--.x.--.x
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'--'x'--'x'--'x`

const c0 = `    .**//////////********//////////,.
xxx,**//////////////////***////////**.
xx.*,    .*.       .***.     ,,.    ..
xx,,**,.,,*,.......,***,.....,,,,,,*,,
xx,...             .,              ,,,.
xx,,,              .,              .,,.
xx,,,               ,               .,,.
x.*,.               ,                ,,.
x.*.               .,                **.
x,,                .,                ,/,
x..                 ,                .*.
x..                .,                .,.
x..                ..                .,.
x..                .,                .,. 
x..                ,*.               .,.
x..               .,,.               .,.
xx.               .,..               .,.
xx.               .,..               .,.
xx.               ...,               .,
xx,               ..x,.              .,
xx,               ,.x,.              .,
xx,               ,.x*.              ..
xx.               ,.x,.              ..
xx,.              ,.x,,              ..
xx,.              ..x.,.             ,.
xx,.              ,xx.,.             ,.
xx,.              ,xx.,.             ,.
xx,.              .xx.,.             *.
xx,.              .xx.,.             ,.
xx,.             ..xx.,.             ,.
xx,.             ..xx.,.             ,.
xx..             ..xx.,.             .
xx.,             ..xx.,.             .
xx.,.            ..xx.,.             .
xx.,.            ..xx.,.             .
xx.,.            ..xxx..             .
xx.,.            .xxxx..             .
xx.,.            .xxxx..            ..
xx.,.           .,xxxx.,            ..
xx.,.           .,xxxx.,            ..
xx.,.           .,xxxx.,            .
xxx,.           .,xxxx.,            ..
xxx..           .,xxxx.,            ..
xxx..           .,xxxx.,            .
xxx.,           ,,xxxx.,            .
xxx.,           ,,xxxx.,            .
xxx.,...........,,xxxx.,............,
xxxx,///////////*.xxxxx.**/////////*,`

const c2 = `   .,,*/*,,,,,,,,,,**,,,,,,,,,,*/*,,.
xxx.  ,*.          **.         .*.  ,.
xx.*,.,,.....,,,,,,//,,,,,,.....,,,.*,
xx..               ,.               .,
xx,                ,.    .,......... .
x.,                ,.    *.       ., ..
x.,                ,.    *.       ., ..
x.,                ,.    *.       ., ..
x..                ..    *.       .,  .
.,.                ,.    ,,.......,,  .
.,.                ,.                 .
.,.                ..                 .
.,.                ,.                 .
x..              ..*,..               .
x,.                *,.                .
x.,                ,,.               ..
x.,               ...,               ..
x.,               ...,               ..
x..               ...,               ..
x.,               .x.,               ..
x.,               .x.,               ..
x..               .xx.               .
xx.               .xx.              ..
xx.               .xx.               .
xx.              .,xx.               .
xx.               .xx.               .
xx.              .,xx..              .
xx..             .,xx..             .,
xx..             ,,xx..             .,
xx.              ,,xx..             .,
xx..             ,,xx.,             .,
xx..             ,.xxx.             .,
xx..             ,.xxx.             .,
xx..             ,.xxx.             .,
xx..             ,.xxx.             ..
xxx.             ,.xxx,.            .,
xxx.             ,.xxx,.            ..
xxx.             ,.xxx,.            ,.
xxx.             ,.xxx,.            ,.
xxx.             .xxxx,.            ,.
xxx.             .xxxx,.            ,.
xxx,             .xxxx,.            ,.
xxx,.            .xxxx,.            ,.
xxx*.           ..xxxx,.            *.
xxx*.           ..xxxx..            ,.
xxx,.           ..xxxx..            ,.
xxx,.           ..xxxx..            ,.
xxx.,,,,,,,,,,,,,xxxxxx,,,,,,,,,,,,,.`

const c1 = `xxxxxxxxxxxxx..,,,*/*,,,,*,
xxxxxx..,**,       ,,    .**
xxxxxx//,.,*.,***,.,,.....,*
xxxxxx*/,,**. ,...         ,
xxxxxx.,.     .,,.         ..
xxxxxx.,.     .,,.         ..
xxxxxxx.,      .*.         ..
xxxxxxx.,      .*.         ..
xxxxxxx.,       ,,         ,
xxxxxxx.,       ..        .*
xxxxxxx.,        .       .,,
xxxxxxx.,        .       .,.
xxxxxxx.,        .       .,.
xxxxxxx.,        .       .,.
xxxxxxx.,        .       ..
xxxxxxx.,        .       ,.
xxxxxxxx.        .       ,.
xxxxxxxx.        .       ,.
xxxxxxxx.        .       ,.
xxxxxxxx.        .       .
xxxxxxxx.        .      ..
xxxxxxxx..       .      ..
xxxxxxxx..       .      ..
xxxxxxxx..       ,      .
xxxxxxxx..       ,      .
xxxxxxxxx.      .,.     ,
xxxxxxxxx.       ,      .
xxxxxxxxx.       .      .
xxxxxxxxx,.      .      .
xxxxxxxxx,.      .      .
xxxxxxxxx,.      .      .
xxxxxxxxx,.      .      .
xxxxxxxxx..      .      .
xxxxxxxxx..      ,      .
xxxxxxxxx..      ,      .
xxxxxxxxx.,.     .      .
xxxxxxxxx.,.     .      ,
xxxxxxxxx.,.     .      ,
xxxxxxxxx.,.     ,      .
xxxxxxxxx.,.     ,      .
xxxxxxxxxx,.     ,     .,
xxxxxxxxxx,.     ,     .,
xxxxxxxxxx,.     ,     .,
xxxxxxxxxx,,     ,     .,
xxxxxxxxxx..     ,.    .,
xxxxxxxxxx..     ,.    .,
xxxxxxxxxx.,     ,.    .,
xxxxxxxxxxxx...,,*,,,,,,.`

const c3 = `xxxxxxxx,*,,,,*/*,,,..
xxxxxxx**.    ,,       ,**,..
xxxxxxx*,.....,,.,***,.*,.,//
xxxxxxx,         ..., .**,,/*
xxxxxx..         .,,.     .,.
xxxxxx..         .,,.     .,.
xxxxxx..         .*.      ,.
xxxxxx..         .*.      ,.
xxxxxxx,         ,,       ,.
xxxxxxx*.        ..       ,.
xxxxxxx,,.       .        ,.
xxxxxxx.,.       .        ,.
xxxxxxx.,.       .        ,.
xxxxxxx.,.       .        ,.
xxxxxxxx..       .        ,.
xxxxxxxx.,       .        ,.
xxxxxxxx.,       .        .
xxxxxxxx.,       .        .
xxxxxxxx.,       .        .
xxxxxxxxx.       .        .
xxxxxxxxx..      .        .
xxxxxxxxx..      .       ..
xxxxxxxxx..      .       ..
xxxxxxxxxx.      ,       ..
xxxxxxxxxx.      ,       ..
xxxxxxxxxx,     .,.      .
xxxxxxxxxx.      ,       .
xxxxxxxxxx.      .       .
xxxxxxxxxx.      .      .,
xxxxxxxxxx.      .      .,
xxxxxxxxxx.      .      .,
xxxxxxxxxx.      .      .,
xxxxxxxxxx.      .      ..
xxxxxxxxxx.      ,      ..
xxxxxxxxxx.      ,      ..
xxxxxxxxxx.      .     .,.
xxxxxxxxxx,      .     .,.
xxxxxxxxxx,      .     .,.
xxxxxxxxxx.      ,     .,.
xxxxxxxxxx.      ,     .,.
xxxxxxxxxx,.     ,     .,
xxxxxxxxxx,.     ,     .,
xxxxxxxxxx,.     ,     .,
xxxxxxxxxx,.     ,     ,,
xxxxxxxxxx,.    .,     ..
xxxxxxxxxx,.    .,     ..
xxxxxxxxxx,.    .,     ,.
xxxxxxxxxx.,,,,,,*,,...`

type Pants struct {
	sprite.BaseSprite
	Timer   int
	TimeOut int
}

type PantsText struct {
	sprite.BaseSprite
	Timer int
}

type Confetti struct {
	sprite.BaseSprite
	Timer     int
	TimeOut   int
	VY        int
	VYTimer   int
	VYTimeOut int
}

func NewPants() *Pants {
	s := &Pants{BaseSprite: sprite.BaseSprite{
		Visible:        true,
		Y:              2,
		CurrentCostume: 0},
		Timer:   0,
		TimeOut: 6}
	s.AddCostume(sprite.NewCostume(c0, 'x'))
	s.AddCostume(sprite.NewCostume(c3, 'x'))
	s.AddCostume(sprite.NewCostume(c2, 'x'))
	s.AddCostume(sprite.NewCostume(c1, 'x'))
	s.X = Width/2 - 20
	return s
}

func (s *Pants) Update() {
	s.Timer = s.Timer + 1
	if s.Timer > s.TimeOut {
		s.CurrentCostume = s.CurrentCostume + 1
		if s.CurrentCostume >= len(s.Costumes) {
			s.CurrentCostume = 0
		}
		s.Timer = 0
	}
}

func NewPantsText() *PantsText {
	s := &PantsText{BaseSprite: sprite.BaseSprite{
		Alpha:          'x',
		CurrentCostume: 0},
		Timer: 0}
	s.AddCostume(sprite.NewCostume(pants_c0, 'x'))
	s.X = Width/2 - s.Width/2
	s.Y = Height/2 - s.Height/2
	return s
}

func (s *PantsText) Update() {
	s.Timer = s.Timer + 1
	if s.Timer > 3 {
		s.Visible = !s.Visible
		s.Timer = 0
	}
}

func NewConfetti() *Confetti {
	s := &Confetti{BaseSprite: sprite.BaseSprite{
		Visible: true},
		Timer:   0,
		TimeOut: 3}
	s.AddCostume(sprite.NewCostume(confetti_c0, 'x'))
	s.AddCostume(sprite.NewCostume(confetti_c1, 'x'))
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	s.X = rnd.Intn(Width)
	s.Y = -rnd.Intn(Height)
	s.VY = rnd.Intn(2) + 1
	s.VYTimer = 0
	s.VYTimeOut = 2
	return s
}

func (s *Confetti) Update() {
	s.Timer = s.Timer + 1
	s.VYTimer = s.VYTimer + 1
	if s.Timer > 2 {
		if s.CurrentCostume == 0 {
			s.CurrentCostume = 1
		} else {
			s.CurrentCostume = 0
		}
		s.Timer = 0
	}
	if s.VYTimer > s.VYTimeOut {
		s.Y = s.Y + s.VY
		if s.Y > Height {
			s.Y = 0 - s.Height
		}
		s.VYTimer = 0
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

	Width, Height = tm.Size()

	event_queue := make(chan tm.Event)
	go func() {
		for {
			event_queue <- tm.PollEvent()
		}
	}()

	for n := 0; n < 40; n++ {
		c := NewConfetti()
		allSprites.Sprites = append(allSprites.Sprites, c)
	}
	p := NewPants()
	pt := NewPantsText()
	allSprites.Sprites = append(allSprites.Sprites, p)
	allSprites.Sprites = append(allSprites.Sprites, pt)

mainloop:
	for {
		tm.Clear(tm.ColorDefault, tm.ColorDefault)

		select {
		case ev := <-event_queue:
			if ev.Type == tm.EventKey {
				if ev.Key == tm.KeyEsc || ev.Ch == 'q' {
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
