package main

import (
	"math"
	"os"
	"sort"
	"time"

	sprite "github.com/pdevine/go-asciisprite"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

var allSprites sprite.SpriteGroup
var Width int
var Height int

// rotX/rotY are the shared scene rotation angles; every object spins together.
var rotX, rotY float64
// dist is the camera distance from the scene centre (adjustable with +/-).
var dist float64 = 120

type PicoCADDemo struct {
	sprite.BaseSprite
	model   *sprite.PicoCADModel
	orig    [][3]float64 // normalised model-space verts, unrotated
	points  []*Point3D   // rotated this frame
	// faces sorted each frame
	faces   []*sprite.PicoCADFace
}

type Point3D struct {
	X     float64
	Y     float64
	Z     float64
	cX    float64
	cY    float64
	cZ    float64
	fl    float64
	vpX   float64
	vpY   float64
	scale float64
}

func NewPoint3D(x, y, z float64) *Point3D {
	p := &Point3D{
		fl:    250.0,
		scale: 1.0,
		X:     x,
		Y:     y,
		Z:     z,
	}
	return p
}

func (p *Point3D) SetVanishingPoint(vpX, vpY int) {
	p.vpX = float64(vpX)
	p.vpY = float64(vpY)
}

func (p *Point3D) SetCenter(cX, cY, cZ float64) {
	p.cX = cX
	p.cY = cY
	p.cZ = cZ
}

func (p *Point3D) ScreenX() int {
	p.scale = p.fl / (p.fl + p.Z + p.cZ)
	return int(math.Round(p.vpX + (p.cX+p.X)*p.scale))
}

func (p *Point3D) ScreenY() int {
	p.scale = p.fl / (p.fl + p.Z + p.cZ)
	return int(math.Round(p.vpY + (p.cY+p.Y)*p.scale))
}

func rotateXYZ(v [3]float64, rx, ry float64) [3]float64 {
	c, s := math.Cos(rx), math.Sin(rx)
	y := v[1]*c - v[2]*s
	z := v[2]*c + v[1]*s
	c, s = math.Cos(ry), math.Sin(ry)
	x := v[0]*c - z*s
	z = z*c + v[0]*s
	return [3]float64{x, y, z}
}

// shadeForFace computes a flat shade factor for a triangle from the dot
// product of its (world-space) normal with a directional light. Returns
// a value in the range [0.25, 1.0].
func shadeForFace(a, b, c *Point3D) float64 {
	u := [3]float64{b.X - a.X, b.Y - a.Y, b.Z - a.Z}
	v := [3]float64{c.X - a.X, c.Y - a.Y, c.Z - a.Z}
	n := [3]float64{
		u[1]*v[2] - u[2]*v[1],
		u[2]*v[0] - u[0]*v[2],
		u[0]*v[1] - u[1]*v[0],
	}
	l := math.Sqrt(n[0]*n[0] + n[1]*n[1] + n[2]*n[2])
	if l == 0 {
		return 1.0
	}
	// light direction
	lx, ly, lz := 1/math.Sqrt(6), -1/math.Sqrt(6), 1/math.Sqrt(6)*2
	d := (n[0]*lx + n[1]*ly + n[2]*lz) / l
	if d < 0 {
		d = -d
	}
	return 0.25 + 0.75*d
}

func faceDepth(points []*Point3D, f *sprite.PicoCADFace) float64 {
	z := points[f.Verts[0]].Z
	z = math.Min(z, points[f.Verts[1]].Z)
	z = math.Min(z, points[f.Verts[2]].Z)
	return z
}

func NewPicoCADDemo(fn string, ox, oy float64) *PicoCADDemo {
	model, err := sprite.LoadPicoCAD(fn)
	if err != nil {
		panic(err)
	}

	d := &PicoCADDemo{
		BaseSprite: sprite.BaseSprite{Visible: true},
		model:      model,
	}

	// normalise the model's size to about 45 screen cells and centre it
	var minX, maxX, minY, maxY, minZ, maxZ float64
	for i, v := range model.Verts {
		if i == 0 {
			minX, maxX = v[0], v[0]
			minY, maxY = v[1], v[1]
			minZ, maxZ = v[2], v[2]
			continue
		}
		minX = math.Min(minX, v[0])
		maxX = math.Max(maxX, v[0])
		minY = math.Min(minY, v[1])
		maxY = math.Max(maxY, v[1])
		minZ = math.Min(minZ, v[2])
		maxZ = math.Max(maxZ, v[2])
	}
	size := math.Max(maxX-minX, math.Max(maxY-minY, maxZ-minZ))
	scale := 45.0 / size
	cX := (minX + maxX) / 2
	cY := (minY + maxY) / 2
	cZ := (minZ + maxZ) / 2

	for _, v := range model.Verts {
		d.orig = append(d.orig, [3]float64{(v[0] - cX) * scale, (v[1] - cY) * scale, (v[2] - cZ) * scale})
	}

	// rotated points for the first frame
	d.points = make([]*Point3D, len(d.orig))
	for i := range d.points {
		d.points[i] = NewPoint3D(0, 0, 0)
		d.points[i].SetVanishingPoint(Width/2, Height/2)
		d.points[i].SetCenter(ox, oy, dist)
	}

	d.faces = make([]*sprite.PicoCADFace, len(model.Faces))
	copy(d.faces, model.Faces)

	surf := sprite.NewSurface(Width, Height, false)
	d.BlockCostumes = append(d.BlockCostumes, &surf)

	return d
}

func (d *PicoCADDemo) Update() {
	// follow camera distance changes
	for i := range d.points {
		d.points[i].cZ = dist
	}
	for i, v := range d.orig {
		rv := rotateXYZ(v, rotX, rotY)
		d.points[i].X = rv[0]
		d.points[i].Y = rv[1]
		d.points[i].Z = rv[2]
	}

	sort.Slice(d.faces, func(i, j int) bool {
		return faceDepth(d.points, d.faces[i]) > faceDepth(d.points, d.faces[j])
	})

	surf := sprite.NewSurface(Width, Height, false)
	for _, f := range d.faces {
		a := d.points[f.Verts[0]]
		b := d.points[f.Verts[1]]
		c := d.points[f.Verts[2]]

		// backface culling unless the face is marked NoCull
		if !f.NoCull {
			cax := c.ScreenX() - a.ScreenX()
			cay := c.ScreenY() - a.ScreenY()
			bcx := b.ScreenX() - c.ScreenX()
			bcy := b.ScreenY() - c.ScreenY()
			if cax*bcy > cay*bcx {
				continue
			}
		}

		shade := 1.0
		if d.model.Shading && !f.NoShade {
			shade = shadeForFace(a, b, c)
		}

		if f.NoTex {
			surf.Triangle(
				a.ScreenX(), a.ScreenY(),
				b.ScreenX(), b.ScreenY(),
				c.ScreenX(), c.ScreenY(),
				sprite.ShadeRune(f.Flat, shade), true)
			continue
		}
		surf.TexturedTriangle(
			a.ScreenX(), a.ScreenY(),
			b.ScreenX(), b.ScreenY(),
			c.ScreenX(), c.ScreenY(),
			d.model.Texture,
			f.UVs[0][0], f.UVs[0][1],
			f.UVs[1][0], f.UVs[1][1],
			f.UVs[2][0], f.UVs[2][1],
			shade,
		)
	}
	d.BlockCostumes[0] = &surf
}

func main() {
	// XXX - Wait a bit until the terminal is properly initialized
	time.Sleep(500 * time.Millisecond)

	err := tm.Init()
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	w, h := tm.Size()
	Width = w * 2
	Height = h * 2

	event_queue := make(chan tm.Event)
	go func() {
		for {
			event_queue <- tm.PollEvent()
		}
	}()

	models := []string{"../testdata/pig.picocad2", "../testdata/pirate.picocad2"}
	if len(os.Args) > 1 {
		models = os.Args[1:]
	}

	spread := math.Min(float64(Width)/(float64(len(models))+1), 80)
	for i, fn := range models {
		ox := (float64(i) - float64(len(models)-1)/2) * spread
		d := NewPicoCADDemo(fn, ox, 0)
		allSprites.Sprites = append(allSprites.Sprites, d)
	}

	allSprites.Init(Width, Height, true)
	allSprites.BlockMode = true
	allSprites.Background = tm.Attribute(178)

mainloop:
	for {
		tm.Clear(tm.Attribute(178), tm.Attribute(178))

		select {
		case ev := <-event_queue:
			if ev.Type == tm.EventKey {
				if ev.Key == tm.KeyEsc {
					break mainloop
				}
				// + / - zoom the camera in and out
				if ev.Ch == '+' || ev.Ch == '=' {
					dist = math.Max(dist-15, 10)
				} else if ev.Ch == '-' || ev.Ch == '_' {
					dist = math.Min(dist+15, 500)
				}
			} else if ev.Type == tm.EventResize {
				Width = ev.Width * 2
				Height = ev.Height * 2
				allSprites.Resize(Width, Height)
			}
		default:
			rotX += 0.02
			rotY += 0.05
			allSprites.Update()
			allSprites.Render()
			time.Sleep(50 * time.Millisecond)
		}
	}
}
