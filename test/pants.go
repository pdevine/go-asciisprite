package main

import (
	"time"

	sprite "github.com/pdevine/go-asciisprite"

	tm "github.com/nsf/termbox-go"
)

var allSprites sprite.SpriteGroup
var Width int
var Height int

const pants_c0 = `
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx,---.,---.,---.
,------.xxx,---.xx,--.xx,--.,--------.x,---.xx|   ||   ||   |
|  .--. 'x/  O  \x|  ,'.|  |'--.  .--''   .-'x|  .'|  .'|  .'
|  '--' ||  .-.  ||  |' '  |xxx|  |xxx` + "`" + `.  ` + "`" + `-.x|  |x|  |x|  |x
|  | --'x|  |x|  ||  |x` + "`" + `   |xxx|  |xxx.-'    |` + "`" + `--'x` + "`" + `--'x` + "`" + `--'x
` + "`" + `--'xxxxx` + "`" + `--'x` + "`" + `--'` + "`" + `--'xx` + "`" + `--'xxx` + "`" + `--'xxx` + "`" + `-----' .--.x.--.x.--.x
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'--'x'--'x'--'x`

const c0 = `    .**//////////********//////////,.
   ,**//////////////////***////////**.
  .*,    .*.       .***.     ,,.    ..
  ,,**,.,,*,.......,***,.....,,,,,,*,,
  ,...             .,              ,,,.
  ,,,              .,              .,,.
  ,,,               ,               .,,.
 .*,.               ,                ,,.
 .*.               .,                **.
 ,,                .,                ,/,
 ..                 ,                .*.
 ..                .,                .,.
 ..                ..                .,.
 ..                .,                .,. 
 ..                ,*.               .,.
 ..               .,,.               .,.
  .               .,..               .,.
  .               .,..               .,.
  .               ...,               .,
  ,               .. ,.              .,
  ,               ,. ,.              .,
  ,               ,. *.              ..
  .               ,. ,.              ..
  ,.              ,. ,,              ..
  ,.              .. .,.             ,.
  ,.              ,  .,.             ,.
  ,.              ,  .,.             ,.
  ,.              .  .,.             *.
  ,.              .  .,.             ,.
  ,.             ..  .,.             ,.
  ,.             ..  .,.             ,.
  ..             ..  .,.             .
  .,             ..  .,.             .
  .,.            ..  .,.             .
  .,.            ..  .,.             .
  .,.            ..   ..             .
  .,.            .    ..             .
  .,.            .    ..            ..
  .,.           .,    .,            ..
  .,.           .,    .,            ..
  .,.           .,    .,            .
   ,.           .,    .,            ..
   ..           .,    .,            ..
   ..           .,    .,            .
   .,           ,,    .,            .
   .,           ,,    .,            .
   .,...........,,    .,............,
    ,///////////*.     .**/////////*,`

const c2 = `   .,,*/*,,,,,,,,,,**,,,,,,,,,,*/*,,.
   .  ,*.          **.         .*.  ,.
  .*,.,,.....,,,,,,//,,,,,,.....,,,.*,
  ..               ,.               .,
  ,                ,.    .,......... .
 .,                ,.    *.       ., ..
 .,                ,.    *.       ., ..
 .,                ,.    *.       ., ..
 ..                ..    *.       .,  .
.,.                ,.    ,,.......,,  .
.,.                ,.                 .
.,.                ..                 .
.,.                ,.                 .
 ..              ..*,..               .
 ,.                *,.                .
 .,                ,,.               ..
 .,               ...,               ..
 .,               ...,               ..
 ..               ...,               ..
 .,               . .,               ..
 .,               . .,               ..
 ..               .  .               .
  .               .  .              ..
  .               .  .               .
  .              .,  .               .
  .               .  .               .
  .              .,  ..              .
  ..             .,  ..             .,
  ..             ,,  ..             .,
  .              ,,  ..             .,
  ..             ,,  .,             .,
  ..             ,.   .             .,
  ..             ,.   .             .,
  ..             ,.   .             .,
  ..             ,.   .             ..
   .             ,.   ,.            .,
   .             ,.   ,.            ..
   .             ,.   ,.            ,.
   .             ,.   ,.            ,.
   .             .    ,.            ,.
   .             .    ,.            ,.
   ,             .    ,.            ,.
   ,.            .    ,.            ,.
   *.           ..    ,.            *.
   *.           ..    ..            ,.
   ,.           ..    ..            ,.
   ,.           ..    ..            ,.
   .,,,,,,,,,,,,,      ,,,,,,,,,,,,,.`

const c1 = `             ..,,,*/*,,,,*,
      ..,**,       ,,    .**
      //,.,*.,***,.,,.....,*
      */,,**. ,...         ,
      .,.     .,,.         ..
      .,.     .,,.         ..
       .,      .*.         ..
       .,      .*.         ..
       .,       ,,         ,
       .,       ..        .*
       .,        .       .,,
       .,        .       .,.
       .,        .       .,.
       .,        .       .,.
       .,        .       ..
       .,        .       ,.
        .        .       ,.
        .        .       ,.
        .        .       ,.
        .        .       .
        .        .      ..
        ..       .      ..
        ..       .      ..
        ..       ,      .
        ..       ,      .
         .      .,.     ,
         .       ,      .
         .       .      .
         ,.      .      .
         ,.      .      .
         ,.      .      .
         ,.      .      .
         ..      .      .
         ..      ,      .
         ..      ,      .
         .,.     .      .
         .,.     .      ,
         .,.     .      ,
         .,.     ,      .
         .,.     ,      .
          ,.     ,     .,
          ,.     ,     .,
          ,.     ,     .,
          ,,     ,     .,
          ..     ,.    .,
          ..     ,.    .,
          .,     ,.    .,
            ...,,*,,,,,,.`

const c4 = `        ,*,,,,*/*,,,..           
       **.    ,,       ,**,..      
       *,.....,,.,***,.*,.,//      
       ,         ..., .**,,/*      
      ..         .,,.     .,.      
      ..         .,,.     .,.      
      ..         .*.      ,.       
      ..         .*.      ,.       
       ,         ,,       ,.       
       *.        ..       ,.       
       ,,.       .        ,.       
       .,.       .        ,.       
       .,.       .        ,.       
       .,.       .        ,.       
        ..       .        ,.       
        .,       .        ,.       
        .,       .        .        
        .,       .        .        
        .,       .        .        
         .       .        .        
         ..      .        .        
         ..      .       ..        
         ..      .       ..        
          .      ,       ..        
          .      ,       ..        
          ,     .,.      .         
          .      ,       .         
          .      .       .         
          .      .      .,         
          .      .      .,         
          .      .      .,         
          .      .      .,         
          .      .      ..         
          .      ,      ..         
          .      ,      ..         
          .      .     .,.         
          ,      .     .,.         
          ,      .     .,.         
          .      ,     .,.         
          .      ,     .,.         
          ,.     ,     .,          
          ,.     ,     .,          
          ,.     ,     .,          
          ,.     ,     ,,          
          ,.    .,     ..          
          ,.    .,     ..          
          ,.    .,     ,.          
          .,,,,,,*,,...`

type Pants struct {
	sprite.BaseSprite
	Timer int
}

type PantsText struct {
	sprite.BaseSprite
	Timer int
}

func NewPants() *Pants {
	s := &Pants{BaseSprite: sprite.BaseSprite{
		Alpha:          'x',
		Visible:        true,
		Y:              2,
		CurrentCostume: 0},
		Timer: 0}
	s.AddCostume(sprite.Costume{c0})
	s.AddCostume(sprite.Costume{c1})
	s.AddCostume(sprite.Costume{c2})
	s.AddCostume(sprite.Costume{c4})
	s.X = Width/2 - 20
	return s
}

func (s *Pants) Update() {
	s.Timer = s.Timer + 1
	if s.Timer > 6 {
		s.CurrentCostume = s.CurrentCostume + 1
		if s.CurrentCostume >= 4 {
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
	s.AddCostume(sprite.Costume{pants_c0})
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

func main() {
	// XXX - Wait a bit until the terminal is properly initialized
	time.Sleep(200 * time.Millisecond)

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
