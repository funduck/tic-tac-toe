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

// PromptMove prompts the user for a move and validates the input
// Returns row, col, giveUp flag, and error
func (i *InputService) PromptMove() (row, col int, giveUp bool, err error) {
	if !i.scanner.Scan() {
		return 0, 0, false, errors.New("failed to read input")
	}

	line := strings.TrimSpace(i.scanner.Text())

	// Check for quit/give up
	if line == "q" || line == "quit" {
		return 0, 0, true, nil
	}

	// Can accept input in format "a1" or "1a" for convenience
	if len(line) == 2 {
		letter := line[0]
		digit := line[1]
		if digit >= '0' && digit <= '2' && letter >= 'a' && letter <= 'c' {
			row = int(letter - 'a')
			col = int(digit - '0')
			return row, col, false, nil
		}
		digit, letter = letter, digit
		if digit >= '0' && digit <= '2' && letter >= 'a' && letter <= 'c' {
			row = int(letter - 'a')
			col = int(digit - '0')
			return row, col, false, nil
		}
	}

	return 0, 0, false, errors.New("invalid input format, expected format like 'a1' or '1a'")
}
