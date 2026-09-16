//go:build darwin || linux

package termbox

// Unix terminal layer: /dev/tty, termios raw mode, SIGWINCH resize.

import (
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

type ttyState struct {
	in      *os.File
	out     *os.File
	saved   unix.Termios
	winch   chan os.Signal
	started bool
}

func (t *ttyState) start() error {
	var err error
	if t.in, err = os.OpenFile("/dev/tty", os.O_RDONLY, 0); err != nil {
		return err
	}
	if t.out, err = os.OpenFile("/dev/tty", os.O_WRONLY, 0); err != nil {
		t.in.Close()
		return err
	}
	term, err := unix.IoctlGetTermios(int(t.out.Fd()), ioctlGetTermios)
	if err != nil {
		t.in.Close()
		t.out.Close()
		return err
	}
	t.saved = *term

	raw := *term
	raw.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK |
		unix.ISTRIP | unix.INLCR | unix.IGNCR |
		unix.ICRNL | unix.IXON
	raw.Oflag &^= unix.OPOST
	raw.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	raw.Cflag &^= unix.CSIZE | unix.PARENB
	raw.Cflag |= unix.CS8
	raw.Cc[unix.VMIN] = 1
	raw.Cc[unix.VTIME] = 0
	if err := unix.IoctlSetTermios(int(t.out.Fd()), ioctlSetTermios, &raw); err != nil {
		t.in.Close()
		t.out.Close()
		return err
	}

	t.winch = make(chan os.Signal, 1)
	signal.Notify(t.winch, unix.SIGWINCH)
	t.started = true
	return nil
}

func (t *ttyState) stop() {
	if !t.started {
		return
	}
	t.started = false
	signal.Stop(t.winch)
	close(t.winch)
	unix.IoctlSetTermios(int(t.out.Fd()), ioctlSetTermios, &t.saved)
	t.in.Close()
	t.out.Close()
}

func (t *ttyState) size() (int, int, error) {
	ws, err := unix.IoctlGetWinsize(int(t.out.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}
	return int(ws.Col), int(ws.Row), nil
}
