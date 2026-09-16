//go:build darwin

package termbox

import "golang.org/x/sys/unix"

const ioctlGetTermios = unix.TIOCGETA
const ioctlSetTermios = unix.TIOCSETA
