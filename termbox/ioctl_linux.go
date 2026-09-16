//go:build linux

package termbox

import "golang.org/x/sys/unix"

const ioctlGetTermios = unix.TCGETS
const ioctlSetTermios = unix.TCSETS
