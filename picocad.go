package sprite

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/color"
	"image/png"
	"math"
	"os"
	"strings"

	palette "github.com/pdevine/go-asciisprite/palette"
)

// PicoCADFace is a single triangle of a PicoCAD model.
type PicoCADFace struct {
	Verts   []int        // indices into PicoCADModel.Verts
	UVs     [][2]float64 // per-vertex normalized UV coordinates
	Color   int          // palette index 0-15
	NoTex   bool         // render as a flat colour instead of sampling the texture
	Flat    rune         // rune used for flat-colour faces (valid iff NoTex)
	NoShade bool
	NoCull  bool
}

// PicoCADMotionSegment animates one node property over a time window.
type PicoCADMotionSegment struct {
	Prop     string    // "pos", "rot", or "scale"
	Axises   []string  // "x", "y", "z"
	Start    float64   // seconds
	Stop     float64
	Delta    float64
	Curve    string    // linear, soft, ease in, ease out, instant, ...
	PingPong bool
	Times    int       // repeat count for "sweep"-style segments (0 = none)
}

// PicoCADNode is a node in the scene graph.
type PicoCADNode struct {
	Name      string
	Pos, Rot, Scale [3]float64
	Motions   []*PicoCADMotionSegment
	// local mesh; vert indices are local to this node
	Verts   [][3]float64
	Faces   []*PicoCADFace
	Nodes   []*PicoCADNode // children
}

// PicoCADModel is a parsed PicoCAD model. For v1 files the model is a simple
// flattened list of faces; for v2 the scene graph is preserved in Root so
// animations can be evaluated via TransformsAt.
type PicoCADModel struct {
	Name     string
	Verts    [][3]float64  // flattened (static) verts - v1 path
	Faces    []*PicoCADFace // flattened (static) faces - v1 path
	Root     *PicoCADNode   // v2 scene graph
	Texture  Surface
	Shading  bool
	Duration float64 // seconds, v2 only (0 = not animated)
}

// LoadPicoCAD parses a PicoCAD save file into a PicoCADModel, auto-detecting
// PicoCAD 1 (JSON export with a base64 texture) and PicoCAD 2 (scene-graph
// save with a 128x128 palette-index texture). Faces are triangulated by fan
// around the first vertex, so resulting models contain only triangles. UVs
// are normalized to the 0-1 range.
func LoadPicoCAD(fn string) (*PicoCADModel, error) {
	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	return ParsePicoCAD(data)
}

// ParsePicoCAD parses PicoCAD save data into a PicoCADModel.
func ParsePicoCAD(data []byte) (*PicoCADModel, error) {
	var probe struct {
		Graph json.RawMessage `json:"graph"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("picocad: %w", err)
	}
	if len(probe.Graph) > 0 {
		return parsePicoCAD2(data)
	}
	return parsePicoCAD1(data)
}

// colorRune returns a ColorMap rune for an RGB triplet.
func colorRune(r, g, b float64) rune {
	c := color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 0xff}
	i := palette.NearestIndex(c)
	return getRuneFromColorMap(i)
}

//////////////////////////////////////////////////////////////////////////////
// PicoCAD 1
//////////////////////////////////////////////////////////////////////////////

// picoCADFile is the on-disk JSON structure of a .picoCAD save.
type picoCADFile struct {
	Name    string `json:"name"`
	Shading bool   `json:"shading"`
	Texture struct {
		Data string `json:"data"`
	} `json:"texture"`
	Objects []struct {
		Name string `json:"name"`
		Mesh struct {
			Verts [][3]float64 `json:"v"`
			Faces []struct {
				V  []int        `json:"v"`
				UV [][2]float64 `json:"uv"`
				C  int          `json:"c"`
				NS int          `json:"ns"`
				NC int          `json:"nc"`
			} `json:"f"`
		} `json:"mesh"`
	} `json:"objects"`
}

func parsePicoCAD1(data []byte) (*PicoCADModel, error) {
	var pf picoCADFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("picocad: %w", err)
	}

	// decode the texture
	texData := pf.Texture.Data
	if i := strings.Index(texData, ","); i != -1 {
		texData = texData[i+1:]
	}
	rawPNG, err := base64.StdEncoding.DecodeString(texData)
	if err != nil {
		return nil, fmt.Errorf("picocad: decoding texture: %w", err)
	}
	img, err := png.Decode(bytes.NewReader(rawPNG))
	if err != nil {
		return nil, fmt.Errorf("picocad: decoding texture png: %w", err)
	}

	m := &PicoCADModel{
		Name:    pf.Name,
		Texture: NewSurfaceFromImage(img, false),
		Shading: pf.Shading,
	}

	texW := float64(img.Bounds().Dx())
	texH := float64(img.Bounds().Dy())

	for _, obj := range pf.Objects {
		for _, v := range obj.Mesh.Verts {
			m.Verts = append(m.Verts, v)
		}
		for _, f := range obj.Mesh.Faces {
			// PicoCAD vertex indices are 1-based; convert to 0-based
			vs := make([]int, len(f.V))
			for i, v := range f.V {
				vs[i] = v - 1
			}
			// triangulate the polygon as a fan
			for i := 1; i+1 < len(vs); i++ {
				face := &PicoCADFace{
					Verts:   []int{vs[0], vs[i], vs[i+1]},
					UVs:     [][2]float64{f.UV[0], f.UV[i], f.UV[i+1]},
					Color:   f.C,
					NoShade: f.NS != 0,
					NoCull:  f.NC != 0,
				}
				// normalize UVs from texture pixel space to 0-1
				for j := range face.UVs {
					face.UVs[j][0] = face.UVs[j][0] / texW
					face.UVs[j][1] = face.UVs[j][1] / texH
				}
				m.Faces = append(m.Faces, face)
			}
		}
	}
	return m, nil
}

//////////////////////////////////////////////////////////////////////////////
// PicoCAD 2
//////////////////////////////////////////////////////////////////////////////

// picoCAD2File is the on-disk JSON structure of a PicoCAD 2 save.
type picoCAD2File struct {
	Graph    picoCAD2Node `json:"graph"`
	Metadata struct {
		Name           string  `json:"name"`
		MotionDuration float64 `json:"motion_duration"`
	} `json:"metadata"`
	Texture struct {
		Pixels           string       `json:"pixels"`
		Colors           [][3]float64 `json:"colors"`
		BackgroundColor  int          `json:"background_color"`
		TransparentColor int          `json:"transparent_color"`
	} `json:"texture"`
}

type picoCAD2Node struct {
	Name      string           `json:"name"`
	Visible   bool             `json:"visible"`
	Transform picoCAD2Xform    `json:"transform"`
	Motions   picoCAD2Motions  `json:"motions"`
	Mesh      *picoCAD2Mesh    `json:"mesh"`
	Children  []picoCAD2Node   `json:"children"`
}

type picoCAD2Xform struct {
	Pos   picoCAD2Vec `json:"pos"`
	Rot   picoCAD2Vec `json:"rot"`
	Scale picoCAD2Vec `json:"scale"`
}

type picoCAD2Vec struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// picoCAD2Motions: tracks is 4 lanes of segment lists
type picoCAD2Motions struct {
	Tracks [][]picoCAD2Seg `json:"tracks"`
}

type picoCAD2Seg struct {
	Prop     string   `json:"prop"`
	Axises   []string `json:"axises"`
	Start    float64  `json:"start"`
	Stop     float64  `json:"stop"`
	Delta    float64  `json:"delta"`
	Curve    string   `json:"curve"`
	PingPong bool     `json:"pingpong"`
	Times    int      `json:"times"`
}

type picoCAD2Mesh struct {
	Name     string    `json:"name"`
	Vertices []float64 `json:"vertices"` // flat x,y,z triples
	Faces    []struct {
		VertexIDs []int     `json:"vertex_ids"`
		UVs       []float64 `json:"uvs"` // flat u,v pairs, normalized 0-1
		Color     int       `json:"color"`
		NoTex     bool      `json:"notex"`
		Dbl       bool      `json:"dbl"`
	} `json:"faces"`
}

// transform matrix: [4][4] row-major, last row is [0 0 0 1]
type picoCAD2Mat [4][4]float64

func picoCAD2Identity() picoCAD2Mat {
	var m picoCAD2Mat
	for i := range m {
		m[i][i] = 1
	}
	return m
}

func picoCAD2Mul(a, b picoCAD2Mat) picoCAD2Mat {
	var m picoCAD2Mat
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			var s float64
			for k := 0; k < 4; k++ {
				s += a[i][k] * b[k][j]
			}
			m[i][j] = s
		}
	}
	return m
}

func picoCAD2ScaleMat(s [3]float64) picoCAD2Mat {
	m := picoCAD2Identity()
	m[0][0], m[1][1], m[2][2] = s[0], s[1], s[2]
	return m
}

func picoCAD2RotMat(r [3]float64) picoCAD2Mat {
	// rotate X then Y then Z (matching PicoCAD node order)
	sx, cx := math.Sin(r[0]), math.Cos(r[0])
	sy, cy := math.Sin(r[1]), math.Cos(r[1])
	sz, cz := math.Sin(r[2]), math.Cos(r[2])
	rx := picoCAD2Identity()
	rx[1][1], rx[1][2], rx[2][1], rx[2][2] = cx, -sx, sx, cx
	ry := picoCAD2Identity()
	ry[0][0], ry[0][2], ry[2][0], ry[2][2] = cy, sy, -sy, cy
	rz := picoCAD2Identity()
	rz[0][0], rz[0][1], rz[1][0], rz[1][1] = cz, -sz, sz, cz
	return picoCAD2Mul(rz, picoCAD2Mul(ry, rx))
}

func picoCAD2TransMat(t [3]float64) picoCAD2Mat {
	m := picoCAD2Identity()
	m[0][3], m[1][3], m[2][3] = t[0], t[1], t[2]
	return m
}

func (m picoCAD2Mat) apply(v [3]float64) [3]float64 {
	return [3]float64{
		m[0][0]*v[0] + m[0][1]*v[1] + m[0][2]*v[2] + m[0][3],
		m[1][0]*v[0] + m[1][1]*v[1] + m[1][2]*v[2] + m[1][3],
		m[2][0]*v[0] + m[2][1]*v[1] + m[2][2]*v[2] + m[2][3],
	}
}

// evaluate a curve function at progress p (0-1)
func picoCAD2Curve(curve string, p float64) float64 {
	if p <= 0 {
		return 0
	}
	if p >= 1 {
		return 1
	}
	switch curve {
	case "soft":
		return p * p * (3 - 2*p)
	case "ease in":
		return p * p
	case "ease out":
		return 1 - (1-p)*(1-p)
	case "ease in out":
		if p < 0.5 {
			return 2 * p * p
		}
		return 1 - 2*(1-p)*(1-p)
	case "instant":
		return 1
	default: // "linear" or unknown
		return p
	}
}

// evaluate a list of segments at time t, returning the additive offsets to
// pos/rot/scale (scale offsets are additive toward 1)
func evalMotions(segs []*PicoCADMotionSegment, t float64) (dPos, dRot, dScale [3]float64) {
	// defaults: pos/rot start at 0, scale at 1
	dScale = [3]float64{0, 0, 0}
	for _, s := range segs {
		dur := s.Stop - s.Start
		if dur <= 0 {
			continue
		}
		tt := t - s.Start
		if tt < 0 {
			continue
		}
		// repeat handling: times>0 means the segment repeats times times
		// within its window
		cycles := 1
		if s.Times > 0 {
			cycles = s.Times
		}
		if tt >= dur*float64(cycles) {
			// past the end of all cycles holds at final value only if not
			// pingpong; otherwise fall through and clamp
			if tt >= dur*float64(cycles) {
				tt = dur*float64(cycles) - 1e-9
			}
		}
		// position within one cycle
		c := math.Mod(tt, dur) / dur
		if s.PingPong && int(math.Mod(tt, dur*2)) >= 2 || (s.PingPong && math.Mod(tt, dur*2) >= dur) {
			c = 1 - c
		}
		v := s.Delta * picoCAD2Curve(s.Curve, c)
		for _, ax := range s.Axises {
			i := 0
			switch ax {
			case "y":
				i = 1
			case "z":
				i = 2
			}
			switch s.Prop {
			case "pos":
				dPos[i] += v
			case "rot":
				dRot[i] += v
			case "scale":
				dScale[i] += v
			}
		}
	}
	return
}

// parseNode converts a JSON node into a PicoCADNode tree.
func parseNode2(n picoCAD2Node, runes []rune) *PicoCADNode {
	nd := &PicoCADNode{
		Name:  n.Name,
		Pos:   [3]float64{n.Transform.Pos.X, n.Transform.Pos.Y, n.Transform.Pos.Z},
		Rot:   [3]float64{n.Transform.Rot.X, n.Transform.Rot.Y, n.Transform.Rot.Z},
		Scale: [3]float64{n.Transform.Scale.X, n.Transform.Scale.Y, n.Transform.Scale.Z},
	}
	// picoCAD 2 scale of all-zero means unset; default to 1
	if nd.Scale[0] == 0 && nd.Scale[1] == 0 && nd.Scale[2] == 0 {
		nd.Scale = [3]float64{1, 1, 1}
	}

	for _, track := range n.Motions.Tracks {
		for _, s := range track {
			nd.Motions = append(nd.Motions, &PicoCADMotionSegment{
				Prop:     s.Prop,
				Axises:   s.Axises,
				Start:    s.Start,
				Stop:     s.Stop,
				Delta:    s.Delta,
				Curve:    s.Curve,
				PingPong: s.PingPong,
				Times:    s.Times,
			})
		}
	}

	if n.Visible && n.Mesh != nil {
		for i := 0; i+3 <= len(n.Mesh.Vertices); i += 3 {
			nd.Verts = append(nd.Verts, [3]float64{
				n.Mesh.Vertices[i], n.Mesh.Vertices[i+1], n.Mesh.Vertices[i+2]})
		}
		for _, f := range n.Mesh.Faces {
			vs := make([]int, len(f.VertexIDs))
			for i, id := range f.VertexIDs {
				vs[i] = id - 1 // 1-based -> local 0-based
			}
			uvs := make([][2]float64, len(f.UVs)/2)
			for i := range uvs {
				uvs[i] = [2]float64{f.UVs[i*2], f.UVs[i*2+1]}
			}
			flat := rune(0)
			if f.NoTex && f.Color < len(runes) {
				flat = runes[f.Color]
			}
			for i := 1; i+1 < len(vs); i++ {
				nd.Faces = append(nd.Faces, &PicoCADFace{
					Verts:  []int{vs[0], vs[i], vs[i+1]},
					UVs:    [][2]float64{uvs[0], uvs[i], uvs[i+1]},
					Color:  f.Color,
					NoTex:  f.NoTex,
					Flat:   flat,
					NoCull: f.Dbl,
				})
			}
		}
	}

	for _, c := range n.Children {
		nd.Nodes = append(nd.Nodes, parseNode2(c, runes))
	}
	return nd
}

// Walk applies fn to each node's mesh with the accumulated world matrix.
// fn is called with the node's transformed (world-space) verts and faces.
// Face vert indices are local to the verts slice passed to fn.
// At t=0 the static transforms are applied; for animated models t is the
// time in seconds (wrap with m.Duration).
func (n *PicoCADNode) Walk(fn func(verts [][3]float64, faces []*PicoCADFace), t float64) {
	n.walkMat(picoCAD2Identity(), t, fn)
}

func (n *PicoCADNode) walkMat(parent picoCAD2Mat, t float64, fn func(verts [][3]float64, faces []*PicoCADFace)) {
	dPos, dRot, dScale := evalMotions(n.Motions, t)

	pos := [3]float64{n.Pos[0] + dPos[0], n.Pos[1] + dPos[1], n.Pos[2] + dPos[2]}
	rot := [3]float64{n.Rot[0] + dRot[0], n.Rot[1] + dRot[1], n.Rot[2] + dRot[2]}
	scale := [3]float64{n.Scale[0] + dScale[0], n.Scale[1] + dScale[1], n.Scale[2] + dScale[2]}

	// local matrix: T(pos) * R(rot) * S(scale)
	local := picoCAD2Mul(picoCAD2TransMat(pos),
		picoCAD2Mul(picoCAD2RotMat(rot), picoCAD2ScaleMat(scale)))
	world := picoCAD2Mul(parent, local)

	if len(n.Verts) > 0 {
		wx := make([][3]float64, len(n.Verts))
		for i, v := range n.Verts {
			wx[i] = world.apply(v)
		}
		fn(wx, n.Faces)
	}
	for _, c := range n.Nodes {
		c.walkMat(world, t, fn)
	}
}

func parsePicoCAD2(data []byte) (*PicoCADModel, error) {
	var pf picoCAD2File
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("picocad: %w", err)
	}

	// decode the 128x128 palette-index texture into a Surface
	const texSize = 128
	if len(pf.Texture.Pixels) < texSize*texSize {
		return nil, fmt.Errorf("picocad: texture has %d pixels, expected %d",
			len(pf.Texture.Pixels), texSize*texSize)
	}
	runes := make([]rune, len(pf.Texture.Colors))
	for i, c := range pf.Texture.Colors {
		if i == pf.Texture.TransparentColor {
			continue
		}
		runes[i] = colorRune(c[0], c[1], c[2])
	}

	tex := NewSurface(texSize, texSize, false)
	for i, ch := range pf.Texture.Pixels[:texSize*texSize] {
		var idx int
		fmt.Sscanf(string(ch), "%x", &idx)
		if idx < len(runes) {
			tex.Blocks[i/texSize][i%texSize] = runes[idx]
		}
	}

	m := &PicoCADModel{
		Name:     pf.Metadata.Name,
		Texture:  tex,
		Shading:  true,
		Duration: pf.Metadata.MotionDuration,
		Root:     parseNode2(pf.Graph, runes),
	}

	// keep the flattened view for compatibility: faces index into m.Verts
	m.Root.Walk(func(verts [][3]float64, faces []*PicoCADFace) {
		base := len(m.Verts)
		m.Verts = append(m.Verts, verts...)
		for _, f := range faces {
			nf := *f
			nf.Verts = make([]int, len(f.Verts))
			for i, vi := range f.Verts {
				nf.Verts[i] = base + vi
			}
			m.Faces = append(m.Faces, &nf)
		}
	}, 0)
	return m, nil
}
