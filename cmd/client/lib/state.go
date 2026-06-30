package lib

var UserID string
var Password string
var AccessToken string
var GameID string
var GameState *Game
var Private bool

// LastMove is the 1-based numpad cell of the most recent move, highlighted on
// the board. -1 means no move to highlight yet.
var LastMove = -1
