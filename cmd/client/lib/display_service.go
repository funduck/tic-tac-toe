package lib

import "fmt"

const (
	markEmpty = "."
	markX     = "X"
	markO     = "O"
)

// DisplayService handles rendering game state to the console
type DisplayService struct{}

// NewDisplayService creates a new DisplayService
func NewDisplayService() *DisplayService {
	return &DisplayService{}
}

func markStr(v int32) string {
	switch v {
	case 1:
		return markX
	case 2:
		return markO
	default:
		return markEmpty
	}
}

// MyMark returns the mark (X or O) for the given user
func (d *DisplayService) MyMark() string {
	if GameState.GetUserID1() == UserID {
		return markX
	}
	return markO
}

// PrintBoard renders the game board
func (d *DisplayService) PrintBoard() {
	fmt.Println()
	fmt.Println("  0 1 2")
	board := GameState.GetBoard()
	for r := 0; r < 3; r++ {
		var row []int32
		if r < len(board) {
			row = board[r]
		}
		cells := make([]string, 3)
		for c := 0; c < 3; c++ {
			var v int32
			if c < len(row) {
				v = row[c]
			}
			cells[c] = markStr(v)
		}
		fmt.Printf("%d %s %s %s\n", r, cells[0], cells[1], cells[2])
	}
	fmt.Println()
}

// PrintResult displays the final game result
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

// PrintStatus displays the current game status
func (d *DisplayService) PrintStatus() {
	status := GameState.GetStatus()
	currentPlayerID := GameState.GetCurrentPlayerID()
	switch status {
	case "waiting":
		fmt.Println("⏳ Waiting for opponent to join...")
	case "in_progress":
		if currentPlayerID == UserID {
			fmt.Println("✨ Your turn!")
		} else {
			fmt.Println("⏳ Waiting for opponent's move...")
		}
	case "finished":
		fmt.Println("🏁 Game finished!")
	}
}

// PrintError displays an error message
func (d *DisplayService) PrintError(message string) {
	fmt.Printf("❌ Error: %s\n", message)
}

// PrintInfo displays an info message
func (d *DisplayService) PrintInfo(message string) {
	fmt.Println(message)
}
