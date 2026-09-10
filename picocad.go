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

// PicoCADModel is a parsed PicoCAD model (v1 export or v2 scene-graph save).
type PicoCADModel struct {
	Name    string
	Verts   [][3]float64
	Faces   []*PicoCADFace
	Texture Surface
	Shading bool
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
		Name string `json:"name"`
	} `json:"metadata"`
	Texture struct {
		Pixels           string        `json:"pixels"`
		Colors           [][3]float64  `json:"colors"`
		BackgroundColor  int           `json:"background_color"`
		TransparentColor int           `json:"transparent_color"`
	} `json:"texture"`
}

type picoCAD2Node struct {
	Name      string         `json:"name"`
	Visible   bool           `json:"visible"`
	Transform picoCAD2Xform  `json:"transform"`
	Mesh      *picoCAD2Mesh  `json:"mesh"`
	Children  []picoCAD2Node `json:"children"`
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

func (n picoCAD2Node) apply(v [3]float64) [3]float64 {
	t := n.Transform
	// scale, rotate (x, y, z), then translate
	v[0] *= t.Scale.X
	v[1] *= t.Scale.Y
	v[2] *= t.Scale.Z

	s, c := math.Sin(t.Rot.X), math.Cos(t.Rot.X)
	v[1], v[2] = v[1]*c-v[2]*s, v[2]*c+v[1]*s
	s, c = math.Sin(t.Rot.Y), math.Cos(t.Rot.Y)
	v[0], v[2] = v[0]*c-v[2]*s, v[2]*c+v[0]*s
	s, c = math.Sin(t.Rot.Z), math.Cos(t.Rot.Z)
	v[0], v[1] = v[0]*c-v[1]*s, v[1]*c+v[0]*s

	v[0] += t.Pos.X
	v[1] += t.Pos.Y
	v[2] += t.Pos.Z
	return v
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
		Name:    pf.Metadata.Name,
		Texture: tex,
		Shading: true,
	}

	// walk the scene graph, applying node transforms
	var walk func(n picoCAD2Node, xf func([3]float64) [3]float64)
	walk = func(n picoCAD2Node, xf func([3]float64) [3]float64) {
		if n.Visible {
			apply := func(v [3]float64) [3]float64 {
				return xf(n.apply(v))
			}
			if n.Mesh != nil {
				base := len(m.Verts)
				for i := 0; i+3 <= len(n.Mesh.Vertices); i += 3 {
					v := [3]float64{n.Mesh.Vertices[i], n.Mesh.Vertices[i+1], n.Mesh.Vertices[i+2]}
					m.Verts = append(m.Verts, apply(v))
				}
				for _, f := range n.Mesh.Faces {
					// PicoCAD 2 vertex ids are 1-based; convert to
					// 0-based global indices
					vs := make([]int, len(f.VertexIDs))
					for i, id := range f.VertexIDs {
						vs[i] = base + id - 1
					}
					// uvs are flat u,v pairs
					uvs := make([][2]float64, len(f.UVs)/2)
					for i := range uvs {
						uvs[i] = [2]float64{f.UVs[i*2], f.UVs[i*2+1]}
					}
					flat := rune(0)
					if f.NoTex && f.Color < len(runes) {
						flat = runes[f.Color]
					}
					// triangulate the polygon as a fan
					for i := 1; i+1 < len(vs); i++ {
						m.Faces = append(m.Faces, &PicoCADFace{
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
		}
		for _, c := range n.Children {
			walk(c, xf)
		}
	}
	identity := func(v [3]float64) [3]float64 { return v }
	for _, c := range pf.Graph.Children {
		walk(c, identity)
	}
	return m, nil
}
