package lib

import (
	"bufio"
	"errors"
	"strings"
)

// InputService handles user input
type InputService struct {
	scanner *bufio.Scanner
}

// NewInputService creates a new InputService
func NewInputService(scanner *bufio.Scanner) *InputService {
	return &InputService{
		scanner: scanner,
	}
}

// Lines starts a goroutine that streams input lines from stdin until it closes,
// then closes the returned channel. This lets the game loop react to input
// (e.g. forfeiting with 'q') while it is also polling for opponent moves.
// It is intended to be called once.
func (i *InputService) Lines() <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for i.scanner.Scan() {
			ch <- i.scanner.Text()
		}
	}()
	return ch
}

// PromptMove reads the next line from stdin and interprets it as a move.
// See ParseMove for the accepted formats.
func (i *InputService) PromptMove() (row, col int, giveUp bool, err error) {
	if !i.scanner.Scan() {
		return 0, 0, false, errors.New("failed to read input")
	}
	return ParseMove(i.scanner.Text())
}

// ParseMove interprets a single line of input. Moves use numpad numbering:
// the keys 1-9 map to cells left-to-right, top-to-bottom. "q"/"quit" forfeits.
// Returns row, col (0-based), a giveUp flag, and an error for invalid input.
func ParseMove(line string) (row, col int, giveUp bool, err error) {
	line = strings.TrimSpace(line)

	if line == "q" || line == "quit" {
		return 0, 0, true, nil
	}

	if len(line) == 1 && line[0] >= '1' && line[0] <= '9' {
		n := int(line[0] - '1') // 0-based cell index
		return n / 3, n % 3, false, nil
	}

	return 0, 0, false, errors.New("enter a number 1-9, or q to quit")
}
