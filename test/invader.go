package main

import (
	"fmt"
	"math/rand"
	"time"

	sprite "github.com/pdevine/go-asciisprite"
	tm "github.com/gdamore/tcell/termbox"
)

var allSprites sprite.SpriteGroup
var Width int
var Height int
var gameState *GameState

var randSrc *rand.Rand

type GameState struct {
	invaders  []*Invader
	direction int
	nextFire  *NextInvaderToFire
	player    *Fighter
	lives     []*Fighter
	GameOver  bool
	Score     *Score
	UfoTimer  int
	Ufo       *Ufo
}

type Invader struct {
	sprite.BaseSprite
	Timer     int
	TimeOut   int
	Type	  int
	Col	  int
	Exploding bool
	Dead      bool
}

type NextInvaderToFire struct {
	Invader *Invader
	Timer   int
	TimeOut int
	Fired	bool
}

type Ufo struct {
	sprite.BaseSprite
	Timer   int
	TimeOut int
	Exploding bool
	Dead      bool
}

type Fighter struct {
	sprite.BaseSprite
	Timer     int
	TimeOut   int
	VX        float32
	AX        float32
	Counter   int
	Exploding bool
	Dead      bool
}

type Bullet struct {
	sprite.BaseSprite
	VY	int
	Dead	bool
}

type Score struct {
	sprite.BaseSprite
	Val int
}

/*
const invader_c0 = `  X     X
 XXXXXXXXX
 XX XXX XX
 XXXXXXXXX
 X X   X X
X X   X X`

const invader_c1 = `  X     X
 XXXXXXXXX
 XX XXX XX
 XXXXXXXXX
 X X   X X 
  X X   X X`
*/

const invader_c0 = `  X     X   
   X   X    
  XXXXXXX   
 XX XXX XX  
XXXXXXXXXXX 
X XXXXXXX X 
X X     X X 
   XX XX     `

const invader_c1 = `  X     X   
X  X   X  X 
X XXXXXXX X 
XXX XXX XXX 
XXXXXXXXXXX 
 XXXXXXXXX  
  X     X   
 X       X  `


const invader2_c0 = `   XX     
  XXXX   
 XXXXXX  
XX XX XX 
XXXXXXXX 
  X  X   
 X XX X  
X X  X X  `

const invader2_c1 = `   XX    
  XXXX   
 XXXXXX  
XX XX XX 
XXXXXXXX 
  X  X   
 X    X  
  X  X    `

const invader3_c0 = `    XXXX   
 XXXXXXXXXX  
XXXXXXXXXXXX 
XXX  XX  XXX 
XXXXXXXXXXXX 
   XX  XX    
  X  XX  X   
   X    X    `

const invader3_c1 = `    XXXX   
 XXXXXXXXXX  
XXXXXXXXXXXX 
XXX  XX  XXX 
XXXXXXXXXXXX 
   XX  XX    
  XX XX XX   
XX        XX `


const fighter_c0 = `
      X       
     XXX      
     XXX     
 XXXXXXXXXXX  
XXXXXXXXXXXXX
XXXXXXXXXXXXX
XXXXXXXXXXXXX`


const fighter_explode_c0 = `
   X
     X X X   
   X X       
      XX XX  
X   X XX X   
  XXXXXXXX X 
 XXXXXXXXXX X`

const fighter_explode_c1 = `

             
    X        
  X    X X X 
X   X XX     
  X XX X X X 
 XX XXX XXX X`


const explosion_c1 = `X  X   X  X  
 X  X X  X   
  X     X    
XX       XX  
  X     X    
 X  X X  X   
X  X   X  X 
            `

const explosion_c0 = `        
         
     XX  
   X    X
   X    X
     XX  `


const ufo_c0 = `
     XXXXXXX      
   XXXXXXXXXXX    
  XXXXXXXXXXXXX   
 X XXX XXX XXX X  
XXXXXXXXXXXXXXXXX 
  XXX  XXX  XXX   
   X         X     `

const ufo_c1 = `
     XXXXXXX      
   XXXXXXXXXXX    
  XXXXXXXXXXXXX   
 XXX XXX XXX XXX  
XXXXXXXXXXXXXXXXX 
  XXX  XXX  XXX   
   X         X     `


func NewGame() *GameState {
	gameState = &GameState{player: NewFighter(),}
	gameState.player.X = 34
	allSprites.Sprites = append(allSprites.Sprites, gameState.player)

	gameState.Score = &Score{BaseSprite: sprite.BaseSprite{
		Visible: true,
		X:       20,
		Y:       1},
	}

	gameState.Score.AddCostume(sprite.Convert(sprite.BuildString(fmt.Sprintf("%06d", gameState.Score.Val))))
	allSprites.Sprites = append(allSprites.Sprites, gameState.Score)

	for i := 0; i < 2; i++ {
		f := NewFighter()
		f.Y = 75
		f.X = i*8 + 2
		gameState.lives = append(gameState.lives, f)
		allSprites.Sprites = append(allSprites.Sprites, f)
	}

	return gameState
}

func (gs *GameState) Update() {
	gs.cullSprites()
	gs.checkDirection()
	gs.nextFire.Update()
	if gs.nextFire.Fired {
		gs.setNextInvaderToFire()
	}
	gs.UfoTimer++
	if gs.UfoTimer > 500 {
		gs.UfoTimer = 0
		gs.Ufo = NewUfo()
		allSprites.Sprites = append(allSprites.Sprites, gs.Ufo)
	}
}

func (gs *GameState) setNextInvaderToFire() {
	var invaders map[int]int
	invaders = make(map[int]int)
	cols := []int{}

	// iterate through each of the invaders and determine which columns they are in
	for _, i := range gs.invaders {
		if !i.Dead && !i.Exploding {
			invaders[i.Col] += 1
		}
	}

	// get the keys of the columns and choose one at random
	for k, _ := range invaders {
		cols = append(cols, k)
	}
	if len(cols) == 0 {
		return
	}
	i_col := randSrc.Intn(len(cols))

	// iterate backwards through the invaders and have the first one we find fire
	for i := len(gs.invaders) - 1; i >= 0; i-- {
		if gs.invaders[i].Col == cols[i_col] && !gs.invaders[i].Dead && !gs.invaders[i].Exploding {
			gs.nextFire.Invader = gs.invaders[i]
			break
		}
	}

	gs.nextFire.Timer = 0
	gs.nextFire.TimeOut = 5
	gs.nextFire.Fired = false
}

func (gs *GameState) checkDirection() {
	changeDir := false
	var d int
	for _, s := range gs.invaders {
		if s.X > 98 && gs.direction > 0 {
			d = -(s.X - 98)
			changeDir = true
			break
		} else if s.X < 3  && gs.direction < 0 {
			d = 3 - s.X
			changeDir = true
			break
		}
	}
	if changeDir {
		for _, s := range gs.invaders {
			s.X += d
			if !gs.GameOver {
				s.Y += 3
			}
		}
		gs.direction = -gs.direction
	}
}

func (gs *GameState) cullSprites() {
	for _, s := range allSprites.Sprites {
		switch sp := s.(type) {
		case *Bullet:
			if sp.Dead {
				allSprites.Remove(sp)
			}
		case *Invader:
			if sp.Dead {
				gs.removeInvader(sp)
				allSprites.Remove(sp)
			}
		case *Ufo:
			if sp.Dead {
				allSprites.Remove(sp)
				gs.Ufo = nil
			}
		case *Fighter:
			if sp.Dead {
				allSprites.Remove(sp)
				if len(gs.lives) == 0 {
					gs.GameOver = true
				} else {
					gs.player = gs.lives[len(gs.lives)-1]
					gs.lives = gs.lives[:len(gs.lives)-1]
					gs.player.Y = 70
				}
			}
		default:
		}
	}
}

func (gs *GameState) removeInvader(i *Invader) {
	var idx int
	for cnt, s := range gs.invaders {
		if i == s {
			idx = cnt
			break
		}
	}
	copy(gs.invaders[idx:], gs.invaders[idx+1:])
	gs.invaders[len(gs.invaders)-1] = nil
	gs.invaders = gs.invaders[:len(gs.invaders)-1]
}

func (gs *GameState) createWave() {
	gs.direction = 2

	for cnt := 0; cnt < 9; cnt++ {
		i := NewInvader(0)
		i.Col = cnt
		i.X = cnt * 10 + 10
		i.Y = 7
		gs.invaders = append(gs.invaders, i)
		allSprites.Sprites = append(allSprites.Sprites, i)
	}

	for y := 0; y < 2; y++ {
		for cnt := 0; cnt < 9; cnt++ {
			i := NewInvader(1)
			i.Col = cnt
			i.X = cnt * 10 + 9
			i.Y = y * 6 + 13
			gs.invaders = append(gs.invaders, i)
			allSprites.Sprites = append(allSprites.Sprites, i)
		}
	}

	for y := 0; y < 2; y++ {
		for cnt := 0; cnt < 9; cnt++ {
	  		i := NewInvader(2)
			i.Col = cnt
			i.X = cnt * 10 + 9
			i.Y = y * 6 + 25
			gs.invaders = append(gs.invaders, i)
			allSprites.Sprites = append(allSprites.Sprites, i)
		}
	}

	gs.nextFire = &NextInvaderToFire{
		Invader: gs.invaders[len(gs.invaders)-1],
		TimeOut: 5,
		Fired:   false,
	}
}

func (n *NextInvaderToFire) Update() {
	n.Timer = n.Timer + 1
	if n.Timer > n.TimeOut {
		n.Invader.Fire()
		n.Timer = 0
		n.Fired = true
	}
}

func NewInvader(t int) *Invader {
	s := &Invader{BaseSprite: sprite.BaseSprite{
		Visible:        true,
		X:		2,
		Y:              2,
		CurrentCostume: 0},
		Type:    t,
		Timer:   0,
		TimeOut: 6}

	if t == 1 {
		s.AddCostume(sprite.Convert(invader_c0))
		s.AddCostume(sprite.Convert(invader_c1))
	} else if t == 0 {
		s.AddCostume(sprite.Convert(invader2_c0))
		s.AddCostume(sprite.Convert(invader2_c1))
	} else if t == 2 {
		s.AddCostume(sprite.Convert(invader3_c0))
		s.AddCostume(sprite.Convert(invader3_c1))
	}
	return s
}

func (s *Invader) Explode() {
	s.Costumes = []*sprite.Costume{}
	s.AddCostume(sprite.Convert(explosion_c0))
	s.AddCostume(sprite.Convert(explosion_c1))
	s.Exploding = true
}

func (s *Invader) Update() {
	s.Timer = s.Timer + 1
	if s.Timer > s.TimeOut {
		s.X = s.X + gameState.direction
		s.CurrentCostume = s.CurrentCostume + 1
		if s.CurrentCostume >= len(s.Costumes) {
			s.CurrentCostume = 0
			if s.Exploding {
				s.Dead = true
			}
		}
		s.Timer = 0
	}
}

func (s *Invader) Fire() {
	b := &Bullet{BaseSprite: sprite.BaseSprite{
		Visible: true,
		X: s.X+3,
		Y: s.Y+4},
		VY: 2}
	b.AddCostume(sprite.Convert("X \n X\nX \n  "))
	b.AddCostume(sprite.Convert(" X\nX \n X\n  "))
	allSprites.Sprites = append(allSprites.Sprites, b)
}

func NewUfo() *Ufo {
	s := &Ufo{BaseSprite: sprite.BaseSprite{
		Visible:	true,
		X: -4,
		Y: 4},
		TimeOut: 3}
	s.AddCostume(sprite.Convert(ufo_c0))
	s.AddCostume(sprite.Convert(ufo_c1))
	return s
}

func (s *Ufo) Update() {
	s.X += 1
	s.Timer = s.Timer + 1
	if s.Timer > s.TimeOut {
		s.CurrentCostume = s.CurrentCostume + 1
		if s.CurrentCostume >= len(s.Costumes) {
			s.CurrentCostume = 0
			if s.Exploding {
				s.Dead = true
			}
		}
		s.Timer = 0
	}
}

func (s *Ufo) Explode() {
	s.Costumes = []*sprite.Costume{}
	s.AddCostume(sprite.Convert(explosion_c0))
	s.AddCostume(sprite.Convert(explosion_c1))
	s.Exploding = true
}

func NewFighter() *Fighter {
	s := &Fighter{BaseSprite: sprite.BaseSprite{
		Visible:	true,
		X: 2,
		Y: 70},
		TimeOut: 2,
		}
	s.AddCostume(sprite.Convert(fighter_c0))
	return s
}

func (s *Fighter) MoveLeft() {
	if !s.Exploding || !s.Dead {
		s.AX -= 1
	}
}

func (s *Fighter) MoveRight() {
	if !s.Exploding || !s.Dead {
		s.AX += 1
	}
} 

func (s *Fighter) Fire() {
	if s.Exploding || s.Dead {
		return
	}

	b := &Bullet{BaseSprite: sprite.BaseSprite{
		Visible: true,
		X: s.X+3,
		Y: s.Y},
		VY: -2}
	b.AddCostume(sprite.Convert("X \n  "))
	allSprites.Sprites = append(allSprites.Sprites, b)
}

func (s *Fighter) Update() {
	s.VX = s.VX + s.AX
	s.AX = 0
	s.X += int(s.VX)

	if s.X < 2 {
		s.X = 2
		s.VX = 0
		s.AX = 0
	} else if s.X > 100 {
		s.X = 100
		s.VX = 0
		s.AX = 0
	}

	if s.Exploding {
		s.Timer += 1
		if s.Timer > s.TimeOut {
			s.CurrentCostume = s.CurrentCostume + 1
			if s.CurrentCostume >= len(s.Costumes) {
				s.CurrentCostume = 0
			}
			s.Timer = 0
			s.Counter += 1
			if s.Counter > 4 {
				s.Dead = true
			}
		}
	}
}

func (s *Fighter) Explode() {
	s.Costumes = []*sprite.Costume{}
	s.AddCostume(sprite.Convert(fighter_explode_c0))
	s.AddCostume(sprite.Convert(fighter_explode_c1))
	s.Exploding = true
}

func (s *Bullet) Update() {
	s.Y += s.VY

	if s.Y < 0 || s.Y > 80 {
		s.Dead = true
		return
	}

	s.CurrentCostume = s.CurrentCostume + 1
	if s.CurrentCostume >= len(s.Costumes) {
		s.CurrentCostume = 0
	}

	if gameState.player.HitAtPoint(s.X, s.Y+1) {
		gameState.player.Explode()
	}

	for _, i := range gameState.invaders {
		if !i.Exploding && i.HitAtPoint(s.X, s.Y) {
			if i.Type == 0 {
				gameState.Score.Val += 30
			} else if i.Type == 1 {
				gameState.Score.Val += 20
			} else {
				gameState.Score.Val += 10
			}
			i.Explode()
			s.Dead = true
		}
	}

	if gameState.Ufo != nil && gameState.Ufo.HitAtPoint(s.X, s.Y+1) {
		gameState.Ufo.Explode()
	}

}

func (s *Score) Update() {
	gameState.Score.Costumes = nil
	gameState.Score.AddCostume(sprite.Convert(sprite.BuildString(fmt.Sprintf("%06d", gameState.Score.Val))))
}


func main() {
	// XXX - Wait a bit until the terminal is properly initialized
	time.Sleep(500 * time.Millisecond)

	err := tm.Init()
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	s := rand.NewSource(time.Now().Unix())
	randSrc = rand.New(s)

	Width, Height = tm.Size()

	event_queue := make(chan tm.Event)
	go func() {
		for {
			event_queue <- tm.PollEvent()
		}
	}()


	//ufo := NewUfo()
	//ufo.X = 24
	//allSprites.Sprites = append(allSprites.Sprites, ufo)
	//invaders = append(invaders, ufo)

	gameState := NewGame()


mainloop:
	for {
		tm.Clear(tm.ColorDefault, tm.ColorDefault)

		select {
		case ev := <-event_queue:
			if ev.Type == tm.EventKey {
				if ev.Key == tm.KeyEsc || ev.Ch == 'q' {
					break mainloop
				} else if ev.Key == tm.KeyArrowLeft {
					gameState.player.MoveLeft()
				} else if ev.Key == tm.KeyArrowRight {
					gameState.player.MoveRight()
				} else if ev.Key == tm.KeySpace {
					gameState.player.Fire()
				}
			} else if ev.Type == tm.EventResize {
				Width = ev.Width
				Height = ev.Height
			}
		default:
			if len(gameState.invaders) == 0 {
				gameState.createWave()
			}
			allSprites.Update()
			gameState.Update()
			allSprites.Render()
			time.Sleep(50 * time.Millisecond)
		}
	}
}
