package game

import (
	"errors"
	"slices"
)

const DefaultSize = 3

var (
	ErrCellOccupied = errors.New("cell is occupied")
	ErrOutOfBounds  = errors.New("move out of bounds")
)

// Board is the N×N grid; it owns move validation and win/draw rules.
// 0 = empty, 1 = UserID1 (X), 2 = UserID2 (O).
type Board [][]int

// NewBoard allocates an empty size×size board.
func NewBoard(size int) Board {
	cells := make(Board, size)
	for i := range cells {
		cells[i] = make([]int, size)
	}
	return cells
}

// Size returns the side length of the board.
func (b Board) Size() int {
	return len(b)
}

// inBounds reports whether (x, y) is a valid cell coordinate.
func (b Board) inBounds(x, y int) bool {
	n := b.Size()
	return x >= 0 && x < n && y >= 0 && y < n
}

// isOccupied reports whether the cell at (x, y) already holds a mark.
func (b Board) isOccupied(x, y int) bool {
	return b[x][y] != 0
}

// Set places mark at (x, y).
func (b Board) Set(x, y, mark int) error {
	if mark != 1 && mark != 2 {
		panic("invalid mark: must be 1 or 2")
	}
	if !b.inBounds(x, y) {
		return ErrOutOfBounds
	}
	if b.isOccupied(x, y) {
		return ErrCellOccupied
	}

	b[x][y] = mark
	return nil
}

// HasLine reports whether mark fills a full row, column, or diagonal.
func (b Board) HasLine(mark int) bool {
	n := b.Size()
	if n == 0 {
		return false
	}

	diag, antiDiag := true, true
	for i := range n {
		row, col := true, true
		for j := 0; j < n; j++ {
			if b[i][j] != mark {
				row = false
			}
			if b[j][i] != mark {
				col = false
			}
		}
		if row || col {
			return true
		}
		if b[i][i] != mark {
			diag = false
		}
		if b[i][n-1-i] != mark {
			antiDiag = false
		}
	}
	return diag || antiDiag
}

// IsFull reports whether every cell holds a mark (no empty cells remain).
func (b Board) IsFull() bool {
	for _, row := range b {
		if slices.Contains(row, 0) {
			return false
		}
	}
	return true
}

// IsEmpty reports whether every cell is empty.
func (b Board) IsEmpty() bool {
	for _, row := range b {
		if !slices.Contains(row, 0) {
			return false
		}
	}
	return true
}

// Clone returns a deep copy whose backing rows are independent of b's.
func (b Board) Clone() Board {
	if b == nil {
		return nil
	}
	cells := make(Board, len(b))
	for i, row := range b {
		cells[i] = make([]int, len(row))
		copy(cells[i], row)
	}
	return cells
}
