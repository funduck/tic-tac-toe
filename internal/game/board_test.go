package game

import (
	"errors"
	"testing"
)

func TestNewBoard(t *testing.T) {
	b := NewBoard(3)
	if b.Size() != 3 {
		t.Fatalf("expected size 3, got %d", b.Size())
	}
	for i := range b {
		if len(b[i]) != 3 {
			t.Fatalf("expected row %d to have width 3, got %d", i, len(b[i]))
		}
	}
	if !b.IsEmpty() {
		t.Error("expected a fresh board to be empty")
	}
}

func TestBoardInBounds(t *testing.T) {
	b := NewBoard(3)
	cases := []struct {
		x, y int
		want bool
	}{
		{0, 0, true},
		{2, 2, true},
		{-1, 0, false},
		{0, -1, false},
		{3, 0, false},
		{0, 3, false},
	}
	for _, c := range cases {
		if got := b.inBounds(c.x, c.y); got != c.want {
			t.Errorf("InBounds(%d,%d)=%v, want %v", c.x, c.y, got, c.want)
		}
	}
}

func TestBoardIsOccupied(t *testing.T) {
	b := NewBoard(3)
	if b.isOccupied(1, 1) {
		t.Error("expected empty cell to be unoccupied")
	}
	b.Set(1, 1, 1)
	if !b.isOccupied(1, 1) {
		t.Error("expected cell to be occupied after Set")
	}
}

func TestBoardHasLine(t *testing.T) {
	const mark = 1
	tests := map[string][][2]int{
		"row0":     {{0, 0}, {0, 1}, {0, 2}},
		"row1":     {{1, 0}, {1, 1}, {1, 2}},
		"row2":     {{2, 0}, {2, 1}, {2, 2}},
		"col0":     {{0, 0}, {1, 0}, {2, 0}},
		"col1":     {{0, 1}, {1, 1}, {2, 1}},
		"col2":     {{0, 2}, {1, 2}, {2, 2}},
		"diag":     {{0, 0}, {1, 1}, {2, 2}},
		"antidiag": {{0, 2}, {1, 1}, {2, 0}},
	}
	for name, cells := range tests {
		t.Run(name, func(t *testing.T) {
			b := NewBoard(3)
			for _, c := range cells {
				b.Set(c[0], c[1], mark)
			}
			if !b.HasLine(mark) {
				t.Errorf("expected HasLine to detect %s", name)
			}
			if b.HasLine(2) {
				t.Error("did not expect a line for the other mark")
			}
		})
	}
}

func TestBoardHasLineEmpty(t *testing.T) {
	b := NewBoard(3)
	if b.HasLine(1) {
		t.Error("empty board should not report a line")
	}
}

func TestBoardIsFull(t *testing.T) {
	b := NewBoard(2)
	if b.IsFull() {
		t.Error("expected empty board to not be full")
	}
	for x := 0; x < 2; x++ {
		for y := 0; y < 2; y++ {
			b.Set(x, y, 1)
		}
	}
	if !b.IsFull() {
		t.Error("expected board to be full")
	}
}

func TestBoardCloneIsIndependent(t *testing.T) {
	b := NewBoard(3)
	b.Set(0, 0, 1)

	clone := b.Clone()
	clone.Set(0, 0, 2)
	clone.Set(2, 2, 2)

	if b[0][0] != 1 {
		t.Errorf("mutating clone changed original at (0,0): got %d", b[0][0])
	}
	if b[2][2] != 0 {
		t.Errorf("mutating clone changed original at (2,2): got %d", b[2][2])
	}
}

func TestBoardCloneNil(t *testing.T) {
	var b Board
	if b.Clone() != nil {
		t.Error("expected Clone of nil board to be nil")
	}
}

func TestMakeMove_BoardValidation(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*Game)
		userID  string
		x, y    int
		wantErr error
	}{
		{
			name:   "out of bounds negative x",
			setup:  func(g *Game) {},
			userID: "alice", x: -1, y: 0,
			wantErr: ErrOutOfBounds,
		},
		{
			name:   "out of bounds x > 2",
			setup:  func(g *Game) {},
			userID: "alice", x: 3, y: 0,
			wantErr: ErrOutOfBounds,
		},
		{
			name:   "out of bounds y > 2",
			setup:  func(g *Game) {},
			userID: "alice", x: 0, y: 3,
			wantErr: ErrOutOfBounds,
		},
		{
			name:   "cell occupied",
			setup:  func(g *Game) { g.Board[1][1] = 1 },
			userID: "alice", x: 1, y: 1,
			wantErr: ErrCellOccupied,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGameInProgress("id1", "alice", "bob")
			tt.setup(g)
			err := g.MakeMove(tt.userID, tt.x, tt.y)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}
