package main

import (
	sprite "github.com/pdevine/go-asciisprite"
)

const goomba_walk1 = `      OOOO
     OOOOOO
    OOOOOOOO
   OOOOOOOOOO
  ONNOOOOOONNO
 OOOwNOOOONwOOO
 OOOwNNNNNNwOOO
OOOOwNwOOwNwOOOO
OOOOwwwOOwwwOOOO
OOOOOOOOOOOOOOOO
 OOOOwwwwwwOOOO
    wwwwwwww
    wwwwwwwwNN
   NNwwwwwNNNNN
   NNNwwwNNNNNN
    NNNwwNNNNN`

const goomba_walk2 = `      OOOO
     OOOOOO
    OOOOOOOO
   OOOOOOOOOO
  ONNOOOOOONNO
 OOOwNOOOONwOOO
 OOOwNNNNNNwOOO
OOOOwNwOOwNwOOOO
OOOOwwwOOwwwOOOO
OOOOOOOOOOOOOOOO
 OOOOwwwwwwOOOO
    wwwwwwww
  NNwwwwwwww  
 NNNNNwwwwwNN
 NNNNNNwwwNNN
  NNNNNwwNNN`

const koopa_walk1 = `








   w
  www
  wwwo
 oGwwoo
 oGwwoo
 oGwwoo
 owwwoo
ooowooo
oGooooo
oooooo  GGGGG
ooo oo GoGGGoG
oo  oo GGoGoGGG
oo oowGGGGoGwwG
 o oowGGGoGoGwG
   oowoGoGGGoGo
  oowwGoGGGGGoG
   owGoGoGGGoGo
   owoGGGoGoGGG
    wGGGGGoGGGG
    wwGGGoGoGwww
   oowwwoGGwww
  ooooowwwwwooo
 ooooo      oooo`

type Enemy struct {
	sprite.BaseSprite
	TimeOut int
	Timer   int
}

