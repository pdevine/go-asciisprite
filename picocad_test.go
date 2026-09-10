package sprite

import "testing"

func TestParsePicoCAD(t *testing.T) {
	m, err := LoadPicoCAD("test/testdata/texcube.picocad")
	if err != nil {
		t.Fatalf("LoadPicoCAD: %v", err)
	}

	if len(m.Verts) != 8 {
		t.Errorf("expected 8 verts, got %d", len(m.Verts))
	}

	// 6 quad faces triangulated into 12 triangles
	if len(m.Faces) != 12 {
		t.Errorf("expected 12 faces, got %d", len(m.Faces))
	}

	if m.Texture.Width == 0 || m.Texture.Height == 0 {
		t.Fatalf("texture not decoded, got %dx%d", m.Texture.Width, m.Texture.Height)
	}

	// UVs should be normalized into 0-1 and verts 0-based and in range
	for _, f := range m.Faces {
		for _, uv := range f.UVs {
			if uv[0] < 0 || uv[0] > 1 || uv[1] < 0 || uv[1] > 1 {
				t.Errorf("uv out of range: %v", uv)
			}
		}
		for _, vi := range f.Verts {
			if vi < 0 || vi >= len(m.Verts) {
				t.Errorf("vert index out of range: %d", vi)
			}
		}
	}

	// face 0 is the front fan: uv 6..11 of 12 -> 0.5..~0.91 (blue region)
	f := m.Faces[0]
	surf := NewSurface(64, 64, false)
	surf.TexturedTriangle(0, 0, 63, 0, 63, 63, m.Texture,
		f.UVs[0][0], f.UVs[0][1], f.UVs[1][0], f.UVs[1][1], f.UVs[2][0], f.UVs[2][1], 1.0)

	// sample the centre of the triangle - should be non-zero (blue-ish)
	found := 0
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			if surf.Blocks[y][x] != 0 {
				found++
			}
		}
	}
	if found == 0 {
		t.Errorf("TexturedTriangle drew no pixels")
	}
	t.Logf("TexturedTriangle filled %d pixels", found)

	// shade 0.0 draws nothing; darker shades produce darker colours
	s2 := NewSurface(64, 64, false)
	s2.TexturedTriangle(0, 0, 63, 0, 63, 63, m.Texture,
		f.UVs[0][0], f.UVs[0][1], f.UVs[1][0], f.UVs[1][1], f.UVs[2][0], f.UVs[2][1], 0.0)
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			if s2.Blocks[y][x] != 0 {
				t.Errorf("shade 0 should draw nothing")
			}
		}
	}
	if ShadeRune('R', 0.5) == 'R' {
		t.Errorf("ShadeRune('R', 0.5) should map to a darker rune")
	}
	if ShadeRune('X', 1.0) == 0 {
		t.Errorf("ShadeRune at full brightness should preserve the rune")
	}
	if ShadeRune('?', 1.0) != 0 {
		t.Errorf("ShadeRune of an unknown rune should be 0")
	}
}

func TestParsePicoCAD2(t *testing.T) {
	m, err := LoadPicoCAD("test/testdata/pig.picocad2")
	if err != nil {
		t.Fatalf("LoadPicoCAD: %v", err)
	}

	if len(m.Verts) == 0 || len(m.Faces) == 0 {
		t.Fatalf("empty model: %d verts, %d faces", len(m.Verts), len(m.Faces))
	}

	if m.Texture.Width != 128 || m.Texture.Height != 128 {
		t.Errorf("expected 128x128 texture, got %dx%d", m.Texture.Width, m.Texture.Height)
	}

	// the pig's texture is opaque, so all texture cells should be set
	empty := 0
	for y := 0; y < 128; y++ {
		for x := 0; x < 128; x++ {
			if m.Texture.Blocks[y][x] == 0 {
				empty++
			}
		}
	}
	if empty > 0 {
		t.Errorf("%d transparent texture cells", empty)
	}

	var nt, flat int
	for _, f := range m.Faces {
		if len(f.Verts) != 3 {
			t.Errorf("face not triangulated: %v", f.Verts)
		}
		for _, vi := range f.Verts {
			if vi < 0 || vi >= len(m.Verts) {
				t.Errorf("vert index out of range: %d", vi)
			}
		}
		for _, uv := range f.UVs {
			// v1 UVs are always in range; v2 UVs may stray slightly
			// outside 0-1 and are clamped when sampling
			if uv[0] < -0.5 || uv[0] > 1.5 || uv[1] < -0.5 || uv[1] > 1.5 {
				t.Errorf("uv out of range: %v", uv)
			}
		}
		if f.NoTex {
			nt++
			if f.Flat == 0 {
				t.Errorf("flat face with no colour")
			}
		} else {
			flat++
		}
	}
	t.Logf("%d faces (%d textured, %d flat)", len(m.Faces), flat, nt)

	// textured faces should rasterize: draw the deepest textured face
	surf := NewSurface(256, 256, false)
	for _, f := range m.Faces {
		if f.NoTex {
			continue
		}
		x0 := int((m.Verts[f.Verts[0]][0]+4) * 25)
		y0 := int((m.Verts[f.Verts[0]][1]+4) * 25)
		x1 := int((m.Verts[f.Verts[1]][0]+4) * 25)
		y1 := int((m.Verts[f.Verts[1]][1]+4) * 25)
		x2 := int((m.Verts[f.Verts[2]][0]+4) * 25)
		y2 := int((m.Verts[f.Verts[2]][1]+4) * 25)
		surf.TexturedTriangle(x0, y0, x1, y1, x2, y2, m.Texture,
			f.UVs[0][0], f.UVs[0][1], f.UVs[1][0], f.UVs[1][1], f.UVs[2][0], f.UVs[2][1], 1.0)
	}
	found := 0
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			if surf.Blocks[y][x] != 0 {
				found++
			}
		}
	}
	if found == 0 {
		t.Errorf("no pixels rasterized from textured faces")
	}
	t.Logf("textured faces filled %d pixels total", found)
}

func TestPicoCAD2Animation(t *testing.T) {
	m, err := LoadPicoCAD("test/testdata/helicopter.picocad2")
	if err != nil {
		t.Fatalf("LoadPicoCAD: %v", err)
	}
	if m.Root == nil {
		t.Fatalf("no scene graph")
	}
	if m.Duration != 6.4 {
		t.Errorf("expected duration 6.4, got %v", m.Duration)
	}

	// find the rotor node and check its z-rot animation spins by 20
	// (full delta) between t=1.5 and t=6.4
	var rotor *PicoCADNode
	var findNode func(n *PicoCADNode)
	findNode = func(n *PicoCADNode) {
		if n.Name == "back rotor" {
			rotor = n
		}
		for _, c := range n.Nodes {
			findNode(c)
		}
	}
	findNode(m.Root)
	if rotor == nil {
		t.Fatalf("no back rotor node")
	}
	if len(rotor.Motions) == 0 {
		t.Fatalf("back rotor has no motions")
	}

	// t=0 vs t=end should give different world verts
	var at, bt [3]float64
	rotor.Walk(func(verts [][3]float64, faces []*PicoCADFace) {
		at = verts[0]
	}, 0.0)
	rotor.Walk(func(verts [][3]float64, faces []*PicoCADFace) {
		bt = verts[0]
	}, 6.3)
	if at == bt {
		t.Errorf("rotor did not move between t=0 and t=6.3 (%v)", at)
	}
	t.Logf("rotor vert[0]: t=0 -> %v, t=6.3 -> %v", at, bt)
}
