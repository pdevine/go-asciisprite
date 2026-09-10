package sprite

import "testing"

// exercise the v2 scene-graph path end to end with both real models plus
// a flat-shaded face check
func TestPicoCAD2Scene(t *testing.T) {
	models := []string{"test/testdata/pig.picocad2", "test/testdata/pirate.picocad2"}
	for _, fn := range models {
		m, err := LoadPicoCAD(fn)
		if err != nil {
			t.Fatalf("%s: %v", fn, err)
		}
		surf := NewSurface(200, 200, false)
		count := 0
		for _, f := range m.Faces {
			p0, p1, p2 := m.Verts[f.Verts[0]], m.Verts[f.Verts[1]], m.Verts[f.Verts[2]]
			x0, y0 := int((p0[0]+4)*18), int((p0[1]+4)*18)
			x1, y1 := int((p1[0]+4)*18), int((p1[1]+4)*18)
			x2, y2 := int((p2[0]+4)*18), int((p2[1]+4)*18)
			if f.NoTex {
				surf.Triangle(x0, y0, x1, y1, x2, y2, f.Flat, true)
			} else {
				surf.TexturedTriangle(x0, y0, x1, y1, x2, y2, m.Texture,
					f.UVs[0][0], f.UVs[0][1], f.UVs[1][0], f.UVs[1][1], f.UVs[2][0], f.UVs[2][1], 0.8)
			}
			count++
		}
		found := 0
		for y := 0; y < 200; y++ {
			for x := 0; x < 200; x++ {
				if surf.Blocks[y][x] != 0 {
					found++
				}
			}
		}
		if found == 0 {
			t.Errorf("%s: scene produced no pixels", fn)
		}
		t.Logf("%s: %d faces, %d pixels on scene surf", fn, count, found)
	}
}
