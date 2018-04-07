package main

import "./sprite"

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

func NewWhale(x, y int, costume sprite.Costume) *Whale {
	s := &Whale{BaseSprite: sprite.BaseSprite{
		X:              x,
		Y:              y,
		Height:         0,
		Width:          0,
		Costumes:       []sprite.Costume{},
		CurrentCostume: 0,
	},
	}
	s.AddCostume(costume)
	return s
}

func (s *Whale) Update() {
	s.X = s.X + s.VX
	s.Y = s.Y + s.VY

	if s.X <= 0 {
		s.VX = 1
	}
	if s.X > Width-s.Width {
		s.VX = -1
	}
	if s.Y > Height-s.Height {
		s.VY = -1
	}
	if s.Y < 0 {
		s.VY = 1
	}
}
