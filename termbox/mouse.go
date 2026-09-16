package termbox

// Mouse: SGR (CSI <b;x;yM/m) and X10 (CSI M bxy) parsing.

import "bytes"

// parseMouse decodes a mouse event from the buffer front.
// (partial, complete)
func parseMouse(buf *bytes.Buffer) (bool, bool) {
	b := buf.Bytes()
	if len(b) < 3 {
		return false, false
	}
	if b[0] != '\x1b' || b[1] != '[' {
		return false, false
	}
	if b[2] == '<' {
		// SGR: ESC [ < btn ; x ; y M|m
		i := 3
		for ; i < len(b); i++ {
			c := b[i]
			if c == 'M' || c == 'm' {
				break
			}
			if (c < '0' || c > '9') && c != ';' {
				return false, false
			}
		}
		if i == len(b) {
			return true, false
		}
		var btn, x, y int
		release := b[i] == 'm'
		body := string(b[3:i])
		var parts [3]int
		np := 0
		val := 0
		gotDigit := false
		for _, c := range []byte(body) {
			if c == ';' {
				np++
				val = 0
				gotDigit = false
				continue
			}
			val = val*10 + int(c-'0')
			gotDigit = true
			if np < 3 {
				parts[np] = val
			}
		}
		if !gotDigit || np != 2 {
			return false, false
		}
		btn, x, y = parts[0], parts[1]-1, parts[2]-1
		for j := 0; j <= i; j++ {
			buf.ReadByte()
		}
		post(buildMouse(x, y, btn, release))
		return false, true
	}
	if b[2] == 'M' && len(b) >= 6 {
		// X10: ESC [ M cb cx cy
		btn := int(b[3]) - 32
		x := int(b[4]) - 33
		y := int(b[5]) - 33
		for j := 0; j < 6; j++ {
			buf.ReadByte()
		}
		post(buildMouse(x, y, btn, btn&3 == 3))
		return false, true
	}
	return false, false
}

// buildMouse maps buttons to the pseudo-keys the API exposes.
func buildMouse(x, y, btn int, release bool) Event {
	var k Key
	switch {
	case release:
		k = MouseRelease
	case btn&0x43&3 == 0 && btn&64 == 0:
		k = MouseLeft
	case btn&0x43&3 == 1 && btn&64 == 0:
		k = MouseMiddle
	case btn&64 != 0:
		k = MouseRelease // wheel: treat as release-ish
	default:
		k = MouseRight
	}
	return Event{Type: EventMouse, Key: k, MouseX: x, MouseY: y}
}
