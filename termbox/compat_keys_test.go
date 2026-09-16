// Numeric compatibility between the native Key space and the tcell
// values the previous implementation exposed (consumers may embed them).

package termbox

import (
	"testing"

	"github.com/gdamore/tcell"
)

func TestKeyCompatWithTcell(t *testing.T) {
	pairs := []struct {
		name string
		ours Key
		old  tcell.Key
	}{
		{"ArrowUp", KeyArrowUp, tcell.KeyUp},
		{"ArrowDown", KeyArrowDown, tcell.KeyDown},
		{"ArrowRight", KeyArrowRight, tcell.KeyRight},
		{"ArrowLeft", KeyArrowLeft, tcell.KeyLeft},
		{"PgUp", KeyPgup, tcell.KeyPgUp},
		{"PgDn", KeyPgdn, tcell.KeyPgDn},
		{"Home", KeyHome, tcell.KeyHome},
		{"End", KeyEnd, tcell.KeyEnd},
		{"Insert", KeyInsert, tcell.KeyInsert},
		{"Delete", KeyDelete, tcell.KeyDelete},
		{"F1", KeyF1, tcell.KeyF1},
		{"F12", KeyF12, tcell.KeyF12},
		{"Backspace", KeyBackspace, tcell.KeyBackspace},
		{"Backspace2", KeyBackspace2, tcell.KeyBackspace2},
		{"Tab", KeyTab, tcell.KeyTab},
		{"Enter", KeyEnter, tcell.KeyEnter},
		{"Esc", KeyEsc, tcell.KeyEscape},
		{"CtrlA", KeyCtrlA, tcell.KeyCtrlA},
		{"CtrlZ", KeyCtrlZ, tcell.KeyCtrlZ},
		{"MouseLeft", MouseLeft, tcell.KeyF63},
		{"MouseRight", MouseRight, tcell.KeyF62},
		{"MouseMiddle", MouseMiddle, tcell.KeyF61},
		{"MouseRelease", MouseRelease, tcell.KeyF64},
	}
	for _, p := range pairs {
		if int16(p.ours) != int16(p.old) {
			t.Errorf("%s: ours=%d tcell=%d", p.name, p.ours, p.old)
		}
	}
}

func TestModCompatWithTcell(t *testing.T) {
	if int16(ModShift) != int16(tcell.ModShift) ||
		int16(ModCtrl) != int16(tcell.ModCtrl) ||
		int16(ModAlt) != int16(tcell.ModAlt) ||
		int16(ModMeta) != int16(tcell.ModMeta) {
		t.Fatal("modifier bit layout diverged from tcell")
	}
}
