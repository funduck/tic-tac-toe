package lib

import (
	"bufio"
	"errors"
	"strconv"
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

	// Parse row and column
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return 0, 0, false, errors.New("invalid input: enter row and column separated by a space (e.g. 1 2)")
	}

	row, err1 := strconv.Atoi(parts[0])
	col, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return 0, 0, false, errors.New("invalid input: row and column must be numbers")
	}

	return row, col, false, nil
}
