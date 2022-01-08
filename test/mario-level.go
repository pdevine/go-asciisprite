package main

import (
	"strings"

	sprite "github.com/pdevine/go-asciisprite"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

const level1 = `
                    123               1223                           123                                            123              1223                           123              1223              F
        123         456     12223     4556                123        456                                 123        456     12223    4556                123        456     12223    4556              f 123
        456                 45556                         456                                            456                45556                        456                45556                      f 456
                       ?                                                          bbbbbbbb   bbb?             ?           bbb    b??b                                                        mm        f
                                                                                                                                                                                            mmm        f
                                                                                                                                                                                           mmmm        f
                                                                                                                                                                                          mmmmm        f    ccc
                 ?   b?b?b                      oO         oO                  b?b              b     bb   ?  ?  ?     b          bb      m  m          mm  m            bb?b            mmmmmm        f    wbW
  0                                     oO      pP  0      pP                                       0                                    mm  mm    0   mmm  mm                          mmmmmmm    0   f   cCCCc
 7*9     h        0           oO        pP      pP 7*9     pP      0                               7*9            0                     mmm  mmm  7*9 mmmm  mmm   0 oO              oO mmmmmmmm   7*9  f   bbDbb  0
7*8*9   H   !@@@#7*9    !@#   pP        pP !@@# pP7*8*9    pP!@@@#7*9    !@#               !@@#   7*8*9     !@@@#7*9                   mmmm@@mmmm7*8*mmmmm  mmmm#7*9pP  !@#         pPmmmmmmmmm  7*8*9 m   bbNbb#7*9    !@#
GGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG  GGGGGGGGGGGGGGG   GGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG  GGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG  GGGGGGGGGGGGGGG   GGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG  GGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG`

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

const cloudshrub_topleft = `







             NNN
            Nggg
           Ngggg
          Nggggg
         NNggggg
        Nggggggg
        Nggggggg
         Ngggggg`

const cloudshrub_topmiddle = `       NNNN
     NggggN
   NNggggggN
  NggggggggN N
  NgggggggggNgN
  NggggggbgggggN
 NgggbbgggbggggN
NgggbggggggggggN
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg
gggggggggggggggg`

const cloudshrub_topright = `







N  N
N NgN
gNggN
ggggN N
gggggNgN
gggggggN
gggggggN
ggggggNN`

const shrub_leftramp = `               N
              NG
             NGG
            NGGG
           NGGGG
          NGGGGG
         NGGGGGG
        NGGGGGGG
       NGGGGGGGG
      NGGGGGGGGG
     NGGGGGGGGGG
    NGGGGGGGGGGG
   NGGGGGGGGGGGG
  NGGGGGGGGGGGGG
 NGGGGGGGGGGGGGG
NGGGGGGGGGGGGGGG`

const shrub_rightramp = `N
GN
GGN
GGGN
GGGGN
GGGGGN
GGGGGGN
GGGGGGGN
GGGGGGGGN
GGGGGGGGGN
GGGGGGGGGGN
GGGGGGGGGGGN
GGGGGGGGGGGGN
GGGGGGGGGGGGGN
GGGGGGGGGGGGGGN
GGGGGGGGGGGGGGGN`

const shrub_top = `












     NNNNNN
  NNNGGGGGGNNN
NNGGGGGGGGGGGGNN`

const shrub_greenblock = `GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG`

const shrub_specks = `GGGGGGGGGGGGGNGG
GGGGGGGGGGGGNNNG
GGGGGGGGGGGGNNNG
GGGGGGGGGGGGNNNG
GGGGGGGGGGGGNNNG
GGGGGGGGGNNGGNGG
GGGGGGGGGNNGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG
GGGGGGGGGGGGGGGG`

const pipe_topleft = `NNNNNNNNNNNNNNNN
Nggggggggggggggg
NGGGGGggggggGGGG
NgggGGggggggGggG
NgggGGggggggGggG
NgggGGggggggGggG
NgggGGggggggGggG
NgggGGggggggGggG
NgggGGggggggGggG
NgggGGggggggGggG
NgggGGggggggGggG
NgggGGggggggGggG
NgggGGggggggGggG
NgggGGggggggGggG
NNNNNNNNNNNNNNNN
  NNNNNNNNNNNNNN`

const pipe_topright = `NNNNNNNNNNNNNNNN
gggggggggggggggN
GGGGGGGGGGGGGGGN
GGGGGGGGGgGgGggN
GGGGGGGGGGgGgggN
GGGGGGGGGgGgGggN
GGGGGGGGGGgGgggN
GGGGGGGGGgGgGggN
GGGGGGGGGGgGgggN
GGGGGGGGGgGgGggN
GGGGGGGGGGgGgggN
GGGGGGGGGgGgGggN
GGGGGGGGGGgGgggN
GGGGGGGGGgGgGggN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNN`

const pipe_left = `  NgggGGgggggGgg
  NgggGGgggggGgg
  NgggGGgggggGgg
  NgggGGgggggGgg
  NgggGGgggggGgg
  NgggGGgggggGgg
  NgggGGgggggGgg
  NgggGGgggggGgg
  NgggGGgggggGgg
  NgggGGgggggGgg
  NgggGGgggggGgg
  NgggGGgggggGgg
  NgggGGgggggGgg
  NgggGGgggggGgg
  NgggGGgggggGgg
  NgggGGgggggGgg`

const pipe_right = `GGGGGGGGgGgggN
GGGGGGGGGgGggN
GGGGGGGGgGgggN
GGGGGGGGGgGggN
GGGGGGGGgGgggN
GGGGGGGGGgGggN
GGGGGGGGgGgggN
GGGGGGGGGgGggN
GGGGGGGGgGgggN
GGGGGGGGGgGggN
GGGGGGGGgGgggN
GGGGGGGGGgGggN
GGGGGGGGgGgggN
GGGGGGGGGgGggN
GGGGGGGGgGgggN
GGGGGGGGGgGggN`

const question_block = ` OOOOOOOOOOOOOO
OooooooooooooooN
OoNooooooooooNoN
OooooOOOOOoooooN
OoooOONNNOOooooN
OoooOONooOONoooN
OoooOONooOONoooN
OooooNNoOOONoooN
OooooooOONNNoooN
OooooooOONoooooN
OoooooooNNoooooN
OooooooOOooooooN
OooooooOONoooooN
OoNoooooNNoooNoN
OooooooooooooooN
NNNNNNNNNNNNNNNN`

const brick_block = `wwwwwwwwwwwwwwww
OOOOOOONOOOOOOON
OOOOOOONOOOOOOON
NNNNNNNNNNNNNNNN
OOONOOOOOOONOOOO
OOONOOOOOOONOOOO
OOONOOOOOOONOOOO
NNNNNNNNNNNNNNNN
OOOOOOONOOOOOOON
OOOOOOONOOOOOOON
OOOOOOONOOOOOOON
NNNNNNNNNNNNNNNN
OOONOOOOOOONOOOO
OOONOOOOOOONOOOO
OOONOOOOOOONOOOO
NNNNNNNNNNNNNNNN`

const used_block = `NNNNNNNNNNNNNNNN
NOOOOOOOOOOOOOON
NONOOOOOOOOOONON
NOOOOOOOOOOOOOON
NOOOOOOOOOOOOOON
NOOOOOOOOOOOOOON
NOOOOOOOOOOOOOON
NOOOOOOOOOOOOOON
NOOOOOOOOOOOOOON
NOOOOOOOOOOOOOON
NOOOOOOOOOOOOOON
NOOOOOOOOOOOOOON
NOOOOOOOOOOOOOON
NONOOOOOOOOOONON
NOOOOOOOOOOOOOON
NNNNNNNNNNNNNNNN`

const ground_block = `OwwwwwwwwNOwwwwO
wOOOOOOOONwOOOON
wOOOOOOOONwOOOON
wOOOOOOOONwOOOON
wOOOOOOOONwNOOON
wOOOOOOOONONNNNO
wOOOOOOOONwwwwwN
wOOOOOOOONwOOOON
wOOOOOOOONwOOOON
wOOOOOOOONwOOOON
NNOOOOOONwOOOOON
wwNNOOOONwOOOOON
wOwwNNNNwOOOOOON
wOOOwwwNwOOOOOON
wOOOOOONwOOOOONN
ONNNNNNOwNNNNNNO`

const metal_block = `OwwwwwwwwwwwwwwN
wOwwwwwwwwwwwwNN
wwOwwwwwwwwwwNNN
wwwOwwwwwwwwNNNN
wwwwOOOOOOOONNNN
wwwwOOOOOOOONNNN
wwwwOOOOOOOONNNN
wwwwOOOOOOOONNNN
wwwwOOOOOOOONNNN
wwwwOOOOOOOONNNN
wwwwOOOOOOOONNNN
wwwwOOOOOOOONNNN
wwwNNNNNNNNNONNN
wwNNNNNNNNNNNONN
wNNNNNNNNNNNNNON
NNNNNNNNNNNNNNNO`

const castle_opening = `OOOOONNNNNNOOOON
OOONNNNNNNNNNOON
OONNNNNNNNNNNNON
NNNNNNNNNNNNNNNN
ONNNNNNNNNNNNNNO
ONNNNNNNNNNNNNNO
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN`

const castle_crenulation_open = `wwww       wwwww
OOOw       wOOOO
OOOw       wOOOO
OOOw       wOOOO
OOOw       wOOOO
OOOw       wOOOO
OOOw       wOOOO
NNNwwwwwwwwwNNNN
OOOOOOONOOOOOOON
OOOOOOONOOOOOOON
OOOOOOONOOOOOOON
NNNNNNNNNNNNNNNN
OOONOOOOOOONOOOO
OOONOOOOOOONOOOO
OOONOOOOOOONOOOO
NNNNNNNNNNNNNNNN`

const castle_crenulation_closed = `wwwwOOONOOOwwwww
OOOwOOONOOOwOOOO
OOOwOOONOOOwOOOO
OOOwNNNNNNNwOOOO
OOOwOOOOOOOwOOOO
OOOwOOOOOOOwOOOO
OOOwOOOOOOOwOOOO
NNNwwwwwwwwwNNNN
OOOOOOONOOOOOOON
OOOOOOONOOOOOOON
OOOOOOONOOOOOOON
NNNNNNNNNNNNNNNN
OOONOOOOOOONOOOO
OOONOOOOOOONOOOO
OOONOOOOOOONOOOO
NNNNNNNNNNNNNNNN`

const castle_blackblock = `NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNN`

const castle_window_right = `OOOOOOONNNNNNNNN
OOOOOOONNNNNNNNN
OOOOOOONNNNNNNNN
NNNNNNNNNNNNNNNN
OOONOOOONNNNNNNN
OOONOOOONNNNNNNN
OOONOOOONNNNNNNN
NNNNNNNNNNNNNNNN
OOOOOOONNNNNNNNN
OOOOOOONNNNNNNNN
OOOOOOONNNNNNNNN
NNNNNNNNNNNNNNNN
OOONOOOONNNNNNNN
OOONOOOONNNNNNNN
OOONOOOONNNNNNNN
NNNNNNNNNNNNNNNN`

const castle_window_left = `NNNNNNNNOOOOOOON
NNNNNNNNOOOOOOON
NNNNNNNNOOOOOOON
NNNNNNNNNNNNNNNN
NNNNNNNNOOONOOOO
NNNNNNNNOOONOOOO
NNNNNNNNOOONOOOO
NNNNNNNNNNNNNNNN
NNNNNNNNOOOOOOON
NNNNNNNNOOOOOOON
NNNNNNNNOOOOOOON
NNNNNNNNNNNNNNNN
NNNNNNNNOOONOOOO
NNNNNNNNOOONOOOO
NNNNNNNNOOONOOOO
NNNNNNNNNNNNNNNN`

type Block struct {
	sprite.BaseSprite
}

type QuestionBlock struct {
	Block
	Timer   int
	TimeOut int
}

func InitQuestionBlock(X, Y int) *QuestionBlock {
	b := &QuestionBlock{Block: Block{BaseSprite: sprite.BaseSprite{
		Visible: true,
		X:       X,
		Y:       Y},
	},
	}
	return b
}

func (b *QuestionBlock) Update() {
	b.Timer++
	if b.Timer > b.TimeOut {
		b.CurrentCostume++
		if b.CurrentCostume >= len(b.Costumes) {
			b.CurrentCostume = 0
		}
		b.Timer = 0
	}
}

type BrickBlock struct {
	Block
}

func ParseLevel(l string, bg tm.Attribute) []*Block {

	allBlocks := []*Block{}

	for rcnt, row := range strings.Split(l, "\n") {
		for ccnt, blk := range row {
			if blk == ' ' {
				continue
			}
			b := &Block{BaseSprite: sprite.BaseSprite{
				Visible: true,
				X:       ccnt * 8,
				Y:       rcnt * 8},
			}
			switch blk {
			case '?':
				b.AddCostume(sprite.ColorConvert(question_block, bg))
			case 'b':
				b.AddCostume(sprite.ColorConvert(brick_block, bg))
			case 'c':
				b.AddCostume(sprite.ColorConvert(castle_crenulation_open, bg))
			case 'C':
				b.AddCostume(sprite.ColorConvert(castle_crenulation_closed, bg))
			case 'D':
				b.AddCostume(sprite.ColorConvert(castle_opening, bg))
			case 'm':
				b.AddCostume(sprite.ColorConvert(metal_block, bg))
			case 'N':
				b.AddCostume(sprite.ColorConvert(castle_blackblock, bg))
			case 'f':
				b.AddCostume(sprite.ColorConvert(flagpole, bg))
			case 'F':
				b.AddCostume(sprite.ColorConvert(flagpole_top, bg))
			case 'G':
				b.AddCostume(sprite.ColorConvert(ground_block, bg))
			case 'h':
				b.AddCostume(sprite.ColorConvert(koopa_walk1, bg))
			case 'H':
				b.AddCostume(sprite.ColorConvert(goomba_walk1, bg))
			case 'w':
				b.AddCostume(sprite.ColorConvert(castle_window_right, bg))
			case 'W':
				b.AddCostume(sprite.ColorConvert(castle_window_left, bg))
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
			case '*':
				b.AddCostume(sprite.ColorConvert(shrub_specks, bg))
			case '9':
				b.AddCostume(sprite.ColorConvert(shrub_rightramp, bg))
			case '0':
				b.AddCostume(sprite.ColorConvert(shrub_top, bg))
			case '!':
				b.AddCostume(sprite.ColorConvert(cloudshrub_topleft, bg))
			case '@':
				b.AddCostume(sprite.ColorConvert(cloudshrub_topmiddle, bg))
			case '#':
				b.AddCostume(sprite.ColorConvert(cloudshrub_topright, bg))
			case 'p':
				b.AddCostume(sprite.ColorConvert(pipe_left, bg))
			case 'P':
				b.AddCostume(sprite.ColorConvert(pipe_right, bg))
			case 'o':
				b.AddCostume(sprite.ColorConvert(pipe_topleft, bg))
			case 'O':
				b.AddCostume(sprite.ColorConvert(pipe_topright, bg))
			}
			allSprites.Sprites = append(allSprites.Sprites, b)
			allBlocks = append(allBlocks, b)
		}
	}
	return allBlocks
}
