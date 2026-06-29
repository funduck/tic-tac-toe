package lib

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// ANSI escape codes used for the interactive board.
const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiDim    = "\033[2m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
)

// IsTTY reports whether stdout is an interactive terminal. When it is not
// (piped output, CI, DEBUG logging) the client falls back to linear output.
func IsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// ClearScreen clears the terminal and moves the cursor to the home position.
func ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

// SaveCursor / RestoreCursor remember and return to the cursor position
// (DECSC/DECRC) so we can refresh part of the screen without disturbing the
// line the user is typing on.
func SaveCursor()    { fmt.Print("\033[s") }
func RestoreCursor() { fmt.Print("\033[u") }

// colorize wraps s in the given ANSI code(s). Callers should only use it when
// rendering to a real TTY.
func colorize(s, code string) string {
	return code + s + ansiReset
}
