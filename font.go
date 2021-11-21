package sprite

import (
	"fmt"
	"strconv"
	"strings"
)

// Font provides monospaced Unicode banner fonts
type Font struct {
	Map    map[rune]string
	Width  int
	Height int
}

// NewFont provides a monospaced banner font based upon a mapping and size
func NewFont(m map[rune]string, w, h int) *Font {
	f := &Font{
		Map:    m,
		Width:  w,
		Height: h,
	}

	return f
}

// BuildString provides a unicode block string from an ASCII string.
func (f *Font) BuildString(s string) string {
	t := make([]string, f.Height)
	formatStr := "%-" + strconv.Itoa(f.Width) + "v"

	for _, c := range s {
		cs, ok := f.Map[c]
		if !ok {
			for cnt := 0; cnt < f.Height; cnt++ {
				t[cnt] += fmt.Sprintf(formatStr, "")
			}
		} else {
			chRows := strings.Split(cs, "\n")

			for cnt := 0; cnt < f.Height; cnt++ {
				var l string
				if cnt >= len(chRows) {
					l = ""
				} else {
					l = chRows[cnt]
				}
				t[cnt] += fmt.Sprintf(formatStr, l)
			}
		}
	}

	return strings.Join(t, "\n")
}
