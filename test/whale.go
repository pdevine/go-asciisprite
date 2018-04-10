package main

import (
	"math/rand"
	"time"

	"./sprite"
)

const whale_c0 = `xxxxxxxxxxxxxxx##xxxxxxxx.xxxxx
xxxxxxxxx##x##x##xxxxxxx==xxxxx
xxxxxx##x##x##x##xxxxxx===xxxxx
xx/""""""""""""""""\___/x===xxx
x{                      /xx===x
xx\______ o          __/xxxxxxx
xxxx\    \        __/xxxxxxxxxx
xxxxx\____\______/xxxxxxxxxxxxx`

type Whale struct {
	sprite.BaseSprite
	VX int
	VY int
}

func randPos() (int, int) {
	offset := 20
	var x, y int
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	x = r.Intn(Width-2*offset) + offset
	y = r.Intn(Height-2*offset) + offset
	return x, y
}

func randVec() (int, int) {
	var x, y int
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	n := r.Intn(2)
	x = n
	if x == 0 {
		x = -1
	}

	n = r.Intn(2)
	y = n
	if y == 0 {
		y = -1
	}
	return x, y
}

func NewWhale() *Whale {
	s := &Whale{BaseSprite: sprite.BaseSprite{
		Alpha:          'x',
		Height:         0,
		Width:          0,
		Costumes:       []sprite.Costume{},
		CurrentCostume: 0,
	},
	}
	s.X, s.Y = randPos()
	s.VX, s.VY = randVec()
	s.AddCostume(sprite.Costume{whale_c0})
	return s
}

func (s *Whale) Update() {
	s.X = s.X + s.VX
	s.Y = s.Y + s.VY

	if s.X < 0 {
		s.VX = 1
	}
	if s.X > Width-s.Width {
		s.VX = -1
	}
	if s.Y >= Height-s.Height {
		s.VY = -1
	}
	if s.Y <= 0 {
		s.VY = 1
	}
}
