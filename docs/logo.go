package main

import (
	"math"
	"math/rand"
	"time"

	sprite "github.com/pdevine/go-asciisprite"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

var allSprites sprite.SpriteGroup

var Width int
var Height int

var dust []*Dust
var random *rand.Rand

const fairyCost0 = `
{\xxxx/}
x>`+"`"+`()'<
{@\/|\ @}
x`+"`"+`~'|/~'
xxx//
xxx\\
xxx''`

const fairyCost1 = `
 {\xxxx/}
xx>`+"`"+`()'<
{@ /|\/@}
x`+"`"+`~\|`+"`"+`~'
xxxx\\
xxxx//
xxxx`+"``"

const logoCost = `
    .-_'''-.       ,-----.
   '_( )_   \    .'  .-,  '.
  |(_ o _)|  '  / ,-.|  \ _ \
  . (_,_)/___| ;  \  '_ /  | :  _ _    _ _
  |  |  .-----.|  _`+"`"+`,/ \ _/  | ( ' )--( ' )
  '  \  '-   .': (  '\_/ \   ;(_{;}_)(_{;}_)
   \  `+"`"+`-'`+"`"+`   |  \ `+"`"+`"/  \  ) /  (_,_)--(_,_)
    \        /   '. \_/`+"``"+`".'
     `+"`"+`'-...-'      '-----'

      ____       .-'''-.     _______  .-./`+"`"+`) .-./`+"`"+`)
    .'  __ `+"`"+`.   / _     \   /   __  \ \ .-.')\ .-.')
   /   '  \  \ (`+"`"+`' )/`+"`"+`--'  | ,_/  \__)/ `+"`"+`-' \/ `+"`"+`-' \
   |___|  /  |(_ o _).   ,-./  )       `+"`"+`-'`+"`"+`"`+"`"+` `+"`"+`-'`+"`"+`"`+"`"+`
      _.-`+"`"+`   | (_,_). '. \  '_ '`+"`"+`)     .---.  .---.
   .'   _    |.---.  \  : > (_)  )  __ |   |  |   |
   |  _( )_  |\    `+"`"+`-'  |(  .  .-'_/  )|   |  |   |
   \ (_ o _) / \       /  `+"`"+`-'`+"`"+`-'     / |   |  |   |
    '.(_,_).'   `+"`"+`-...-'     `+"`"+`._____.'  '---'  '---'

        .-'''-. .-------. .-------.   .-./`+"`"+`) ,---------.    .-''-.
       / _     \\  _(`+"`"+`)_ \|  _ _   \  \ .-.')\          \ .'_ _   \
      (`+"`"+`' )/`+"`"+`--'| (_ o._)|| ( ' )  |  / `+"`"+`-' \ `+"`"+`--.  ,---'/ ( `+"`"+` )   '
      _ o _).   |  (_,_) /|(_ o _) /   `+"`"+`-'`+"`"+`"`+"`"+`    |   \  . (_ o _)  |
      (_,_). '. |   '-.-' | (_,_).' __ .---.     :_ _:  |  (_,_)___|
     ,---.  \  :|   |     |  |\ \  |  ||   |     (_I_)  '  \   .---.
     \    `+"`"+`-'  ||   |     |  | \ `+"`"+`'   /|   |    (_(=)_)  \  `+"`"+`-'    /
      \       / /   )     |  |  \    / |   |     (_I_)    \       /
       `+"`"+`-...-'  `+"`"+`---'     ''-'   `+"`"+`'-'  '---'     '---'     `+"`"+`'-..-'
`


type Logo struct {
	sprite.BaseSprite
}

func NewLogo() *Logo {
	l := &Logo{BaseSprite: sprite.BaseSprite{
		Visible: true},
	}
	l.AddCostume(sprite.NewCostume(logoCost, 'x'))
	l.X = Width/2 - l.Width/2
	l.Y = Height/2 - l.Height/2

	return l
}

type Fairy struct {
	sprite.BaseSprite
	AngleX  float64
	AngleY  float64
	CenterX int
	CenterY int
	Range   float64
	VelX    float64
	VelY    float64
	TimeOut int
	Timer   int
}

type Dust struct {
	sprite.BaseSprite
	TimeOut int
	Timer   int
	Dead    bool
}

func NewFairy() *Fairy {
	f := &Fairy{BaseSprite: sprite.BaseSprite{
		Visible: true},
		CenterX: Width/2,
		CenterY: Height/2,
		Range:   15,
		VelX:    0.03,
		VelY:    0.02,
		TimeOut: 20,
	}
	f.AddCostume(sprite.NewCostume(fairyCost0, 'x'))
	f.AddCostume(sprite.NewCostume(fairyCost1, 'x'))
	f.CenterX -= f.Width/2
	f.CenterY -= f.Height/2

	return f
}

func (f *Fairy) Update() {
	if f.Timer > f.TimeOut {
		f.Timer = 0
		f.TimeOut = random.Intn(30)+70
		for cnt := 0; cnt <= random.Intn(3); cnt++ {
			d := f.NewDust(cnt-1)
			dust = append(dust, d)
			allSprites.Sprites = append(allSprites.Sprites, d)
		}
	}
	f.Timer++

	xpos := f.CenterX + int(math.Sin(f.AngleX) * f.Range)
	ypos := f.CenterY + int(math.Sin(f.AngleY) * f.Range)

	if xpos < f.X {
		f.CurrentCostume = 0
	} else if xpos > f.X {
		f.CurrentCostume = 1
	}

	f.X = xpos
	f.Y = ypos
	f.AngleX += f.VelX
	f.AngleY += f.VelY
}

func (f *Fairy) NewDust(xOff int) *Dust {
	d := &Dust{BaseSprite: sprite.BaseSprite{
		Y:       f.Y+2,
		Visible: true},
		TimeOut: random.Intn(10)+8,
	}
	if f.CurrentCostume == 0 {
		d.X = f.X+xOff
	} else if f.CurrentCostume == 1 {
		d.X = f.X+7+xOff
	}
	d.AddCostume(sprite.NewCostume("*", '~'))
	d.AddCostume(sprite.NewCostume(",", '~'))
	d.AddCostume(sprite.NewCostume(".", '~'))
	return d
}

func (d *Dust) Update() {
	d.Y++
	d.Timer++
	if d.Timer > d.TimeOut {
		d.CurrentCostume++
		d.TimeOut = random.Intn(5)+3
		if d.CurrentCostume >= len(d.Costumes) {
			d.Visible = false
			d.Dead = true
		}
	}
}

func Vaccuum() {
	for cnt := len(dust)-1; cnt >= 0; cnt-- {
		d := dust[cnt]
		if d.Dead == true {
			allSprites.Remove(d)
			copy(dust[cnt:], dust[cnt+1:])
			dust[len(dust)-1] = nil
			dust = dust[:len(dust)-1]
		}
	}
}

func main() {
	err := tm.Init()
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	Width, Height = tm.Size()

	random = rand.New(rand.NewSource(time.Now().UnixNano()))

	l := NewLogo()
	f := NewFairy()
	allSprites.Sprites = append(allSprites.Sprites, l)
	allSprites.Sprites = append(allSprites.Sprites, f)

	eventQueue := make(chan tm.Event)
	go func() {
		for {
			eventQueue <- tm.PollEvent()
		}
	}()

mainloop:
	for {
		tm.Clear(tm.ColorDefault, tm.ColorDefault)
		select {
		case ev := <-eventQueue:
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
			time.Sleep(60 * time.Millisecond)
		}
		Vaccuum()
	}

}
