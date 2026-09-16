//go:build darwin || linux

package termbox

import (
	"fmt"
	"os"
	"time"
)

// inputReadLoop reads chunks from the tty, parses them, posts events.
// Returns when quit is closed (file close unblocks the read).
//
// The read runs on its own goroutine and posts chunks; the select loop
// multiplexes chunks against a 50ms keytimer used to flush pending
// partial input (lone Esc, split escape sequences, partial UTF-8).
// Read deadlines on /dev/tty are unreliable in practice, which is why
// this doesn't use SetReadDeadline.
func inputReadLoop(f *os.File, quit chan struct{}) error {
	chunks := make(chan []byte, 16)
	go func() {
		buf := make([]byte, 128)
		for {
			n, err := f.Read(buf)
			if n > 0 {
				c := make([]byte, n)
				copy(c, buf[:n])
				select {
				case chunks <- c:
				case <-quit:
					return
				}
			}
			if err != nil {
				select {
				case chunks <- nil:
				case <-quit:
				}
				return
			}
		}
	}()

	keytimer := time.NewTicker(50 * time.Millisecond)
	defer keytimer.Stop()

	for {
		select {
		case <-quit:
			return nil
		case c, ok := <-chunks:
			if !ok || c == nil {
				return nil // read error / closed
			}
			if debugLog != nil {
				fmt.Fprintf(debugLog, "read %q\n", c)
				debugLog.Sync()
			}
			parseBytes(c)
		case <-keytimer.C:
			expireInput()
		}
	}
}
