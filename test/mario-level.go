package main

import (
	"strings"

	sprite "github.com/pdevine/go-asciisprite"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

const level1 = `
                    123
        123         456     12223
        456                 45556
                       ?



                 ?   b?b?b

 789
78889       12223789    123
GGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG`

const cloud_topleft = `







             NNN
            Nwww
           Nwwww
          Nwwwww
         NNwwwww
        Nwwwwwww
        Nwwwwwww
         Nwwwwww`

const cloud_bottomleft = `          Nwwbww
           Nwwbw
            Nwwb
            Nwww
             NNN`

const cloud_topmiddle = `       NNNN
     NwwwwN
   NNwwwwwwN
  NwwwwwwwwN N
  NwwwwwwwwwNwN
  NwwwwwwbwwwwwN
 NwwwbbwwwbwwwwN
NwwwbwwwwwwwwwwN 
wwwwwwwwwwwwwwww
wwwwwwwwwwwwwwww
wwwwwwwwwwwwwwww
wwwwwwwwwwwwwwww
wwwwwwwwwwwwwwww
wwwwwwwwwwwwwwww
wwwwwwwwwwwwwwww
wwwwwwwwwwwwwwww`

const cloud_bottommiddle = `wwwwwwwwwbwwwwww
wbwwwwwwbwwwwwww
bbbwwwbbbbwwwwbw
wwbbbbbbwbbbbbww
wwwwbbwwwwbbbwww
NwwwwwwNwwwwwwww
 NNwwwN NNwwwwNN
   NNN    NNNN`

const cloud_topright = `







N  N
N NwN
wNwwN
wwwwN N
wwwwwNwN
wwwwwwwN
wwwwwwwN
wwwwwwNN`

const cloud_bottomright = `wwwwwN
wwwwwwN
wwwwwwwN
wwwwwwN
wwwwwNN
NwwNN
 NN`

const shrub_leftramp = `               N
              Ng
             Ngg
            Nggg
           Ngggg
          Nggggg
         Ngggggg
        Nggggggg
       Ngggggggg
      Nggggggggg
     Ngggggggggg
    Nggggggggggg
   Ngggggggggggg
  Nggggggggggggg
 Ngggggggggggggg
Nggggggggggggggg`

const shrub_rightramp = `N
gN
ggN
gggN
ggggN
gggggN
ggggggN
gggggggN
ggggggggN
gggggggggN
ggggggggggN
gggggggggggN
ggggggggggggN
gggggggggggggN
ggggggggggggggN
gggggggggggggggN`

const shrub_greenblock = `gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg`

func ParseLevel(l string, bg tm.Attribute) {

	for rcnt, row := range strings.Split(l, "\n") {
		for ccnt, blk := range row {
			if blk == ' ' {
				continue
			}
			b := &Block{BaseSprite: sprite.BaseSprite{
				Visible: true,
				X:       ccnt*8,
				Y:       rcnt*8},
			}
			switch blk {
			case '?':
				b.AddCostume(sprite.ColorConvert(question_block, bg))
			case 'b':
				b.AddCostume(sprite.ColorConvert(brick_block, bg))
			case 'G':
				b.AddCostume(sprite.ColorConvert(ground_block, bg))
			case '1':
				b.AddCostume(sprite.ColorConvert(cloud_topleft, bg))
			case '2':
				b.AddCostume(sprite.ColorConvert(cloud_topmiddle, bg))
			case '3':
				b.AddCostume(sprite.ColorConvert(cloud_topright, bg))
			case '4':
				b.AddCostume(sprite.ColorConvert(cloud_bottomleft, bg))
			case '5':
				b.AddCostume(sprite.ColorConvert(cloud_bottommiddle, bg))
			case '6':
				b.AddCostume(sprite.ColorConvert(cloud_bottomright, bg))
			case '7':
				b.AddCostume(sprite.ColorConvert(shrub_leftramp, bg))
			case '8':
				b.AddCostume(sprite.ColorConvert(shrub_greenblock, bg))
			case '9':
				b.AddCostume(sprite.ColorConvert(shrub_rightramp, bg))
			}
			allSprites.Sprites = append(allSprites.Sprites, b)
		}
	}
}
