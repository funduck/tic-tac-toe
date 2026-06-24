package lib

import (
	"fmt"
	"os"
	"strings"
	"time"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID"
)

const (
	markX = "X"
	markO = "O"
)

// winLines holds the eight winning triples as 0-based cell indices.
var winLines = [8][3]int{
	{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
	{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
	{0, 4, 8}, {2, 4, 6},
}

// DisplayService handles rendering game state to the console.
type DisplayService struct {
	// redraw is true when stdout is an interactive TTY: frames are cleared and
	// redrawn in place with color. When false (DEBUG or piped output) the client
	// prints linearly so logs stay readable.
	redraw bool
}

// NewDisplayService creates a new DisplayService.
func NewDisplayService() *DisplayService {
	return &DisplayService{
		redraw: IsTTY() && os.Getenv("DEBUG") == "",
	}
}

// MyMark returns the mark (X or O) for the current user.
func (d *DisplayService) MyMark() string {
	if GameState.GetUserID1() == UserID {
		return markX
	}
	return markO
}

// RenderFrame draws a single stable frame: a status header followed by the
// board. On a TTY it clears the screen first so the board updates in place.
func (d *DisplayService) RenderFrame() {
	if d.redraw {
		ClearScreen()
	}
	fmt.Println(d.statusHeader())
	fmt.Print(d.renderBoard())
}

// statusHeader is the persistent line above the board: who you are, the
// opponent, and whose turn it is.
func (d *DisplayService) statusHeader() string {
	opponent := GameState.GetOpponentID()
	if opponent == "" {
		opponent = "—"
	}

	var turn string
	switch GameState.GetStatus() {
	case openapi.StatusWaiting:
		turn = "⏳ Waiting for opponent to join..."
	case openapi.StatusInProgress:
		if GameState.GetCurrentPlayerID() == UserID {
			turn = "✨ Your turn"
		} else {
			turn = "⏳ Opponent's turn..."
		}
	case openapi.StatusFinished:
		turn = "🏁 Game finished"
	}

	if presence := formatLastSeen(GameState.GetOpponentLastSeen()); presence != "" {
		opponent = fmt.Sprintf("%s (%s)", opponent, presence)
	}

	return fmt.Sprintf("You are %s · Opponent: %s\n%s", d.MyMark(), opponent, turn)
}

// formatLastSeen turns an RFC3339 timestamp into a human presence string:
// "online" when seen within a few seconds, otherwise "seen Ns ago" / "seen Nm ago".
// Returns "" when the timestamp is missing or unparseable.
func formatLastSeen(ts *string) string {
	if ts == nil || *ts == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, *ts)
	if err != nil {
		return ""
	}
	d := time.Since(t)
	switch {
	case d < 5*time.Second:
		return "online"
	case d < time.Minute:
		return fmt.Sprintf("seen %ds ago", int(d.Seconds()))
	default:
		return fmt.Sprintf("seen %dm ago", int(d.Minutes()))
	}
}

// renderBoard renders the 3x3 grid using numpad numbering: empty cells show
// their move number (1-9), played cells show X/O. On a TTY the last move and a
// winning line are highlighted.
func (d *DisplayService) renderBoard() string {
	board := GameState.GetBoard()
	win := winningLine(board)

	var b strings.Builder
	b.WriteByte('\n')
	for r := range 3 {
		cells := make([]string, 3)
		for c := range 3 {
			cellNum := r*3 + c + 1
			var v int32
			if r < len(board) && c < len(board[r]) {
				v = board[r][c]
			}
			cells[c] = d.cellString(v, cellNum, win)
		}
		fmt.Fprintf(&b, " %s | %s | %s \n", cells[0], cells[1], cells[2])
		if r < 2 {
			b.WriteString("---+---+---\n")
		}
	}
	b.WriteByte('\n')
	return b.String()
}

// cellString renders a single cell. v is the board value (0 empty, 1 X, 2 O),
// cellNum is the 1-based numpad position, win is the winning line (or nil).
func (d *DisplayService) cellString(v int32, cellNum int, win []int) string {
	switch v {
	case 1, 2:
		mark := markX
		if v == 2 {
			mark = markO
		}
		if d.redraw {
			switch {
			case contains(win, cellNum):
				return colorize(mark, ansiGreen+ansiBold)
			case LastMove == cellNum:
				return colorize(mark, ansiYellow+ansiBold)
			}
		}
		return mark
	default:
		hint := fmt.Sprintf("%d", cellNum)
		if d.redraw {
			return colorize(hint, ansiDim)
		}
		return hint
	}
}

// winningLine returns the 1-based cell numbers of a completed line, or nil.
func winningLine(board [][]int32) []int {
	cell := func(i int) int32 {
		r, c := i/3, i%3
		if r < len(board) && c < len(board[r]) {
			return board[r][c]
		}
		return 0
	}
	for _, ln := range winLines {
		a, b, c := cell(ln[0]), cell(ln[1]), cell(ln[2])
		if a != 0 && a == b && b == c {
			return []int{ln[0] + 1, ln[1] + 1, ln[2] + 1}
		}
	}
	return nil
}

func contains(s []int, v int) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

// PrintResult displays the final game result.
func (d *DisplayService) PrintResult() {
	winnerID := GameState.GetWinnerID()
	switch {
	case winnerID == UserID:
		fmt.Println("🎉 You win!")
	case winnerID != "" && winnerID != UserID:
		fmt.Println("😔 You lose.")
	default:
		fmt.Println("🤝 Draw!")
	}
}

// PrintError displays an error message.
func (d *DisplayService) PrintError(message string) {
	fmt.Printf("❌ Error: %s\n", message)
}

// PrintInfo displays an info message.
func (d *DisplayService) PrintInfo(message string) {
	fmt.Println(message)
}
