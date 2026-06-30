package lib

import "testing"

func TestParseMove(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		wantRow, wantCol int
		wantGiveUp       bool
		wantErr          bool
	}{
		{name: "top-left", input: "1", wantRow: 0, wantCol: 0},
		{name: "top-middle", input: "2", wantRow: 0, wantCol: 1},
		{name: "top-right", input: "3", wantRow: 0, wantCol: 2},
		{name: "center", input: "5", wantRow: 1, wantCol: 1},
		{name: "bottom-right", input: "9", wantRow: 2, wantCol: 2},
		{name: "trims whitespace", input: "  7 ", wantRow: 2, wantCol: 0},
		{name: "quit short", input: "q", wantGiveUp: true},
		{name: "quit long", input: "quit", wantGiveUp: true},
		{name: "zero rejected", input: "0", wantErr: true},
		{name: "out of range", input: "10", wantErr: true},
		{name: "letter rejected", input: "a", wantErr: true},
		{name: "empty rejected", input: "", wantErr: true},
		{name: "old format rejected", input: "a1", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row, col, giveUp, err := ParseMove(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if giveUp != tt.wantGiveUp {
				t.Fatalf("giveUp = %v, want %v", giveUp, tt.wantGiveUp)
			}
			if tt.wantGiveUp {
				return
			}
			if row != tt.wantRow || col != tt.wantCol {
				t.Fatalf("(row,col) = (%d,%d), want (%d,%d)", row, col, tt.wantRow, tt.wantCol)
			}
		})
	}
}
