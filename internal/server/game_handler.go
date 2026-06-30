package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/funduck/tic-tac-toe/internal/game"
)

type GameService interface {
	CreateGame(ctx context.Context, userID string, cmd game.CreateGameCommand) (*game.Game, error)
	JoinGame(ctx context.Context, userID string, cmd game.JoinGameCommand) (*game.Game, error)
	JoinAnyGame(ctx context.Context, userID string) (*game.Game, error)
	GetGame(ctx context.Context, userID string, cmd game.GetGameCommand) (*game.Game, error)
	MakeMove(ctx context.Context, userID string, cmd game.MakeMoveCommand) (*game.Game, error)
	GiveUp(ctx context.Context, userID string, cmd game.GiveUpCommand) (*game.Game, error)
	Quit(ctx context.Context, userID string, cmd game.QuitCommand) (*game.Game, error)
	GetLatestGameForUser(ctx context.Context, userID string) (*game.Game, error)
}

type GameHandler struct {
	svc    GameService
	logger *slog.Logger
}

func NewGameHandler(svc GameService, logger *slog.Logger) *GameHandler {
	return &GameHandler{svc: svc, logger: logger}
}

// CreateGame godoc
//
//	@Summary	Create a new game
//	@ID		createGame
//	@Tags		games
//	@Accept		json
//	@Produce	json
//	@Param		request	body		game.CreateGameCommand	true	"Create game request"
//	@Success	201		{object}	game.Game
//	@Failure	400		{object}	ErrorResponse
//	@Failure	409		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/api/games [post]
func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := userIDFromContext(ctx)
	var cmd game.CreateGameCommand
	if err := parseRequestBody(r, &cmd, h.logger, w); err != nil {
		return // parseRequestBody already wrote the error response
	}

	g, err := h.svc.CreateGame(ctx, userID, cmd)
	if err != nil {
		h.logger.Warn("create game failed", "error", err)
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, g)
}

// JoinGame godoc
//
//	@Summary	Join a waiting game
//	@ID		joinGame
//	@Tags		games
//	@Produce	json
//	@Param		gameID	path		string	true	"Game ID"
//	@Success	200		{object}	game.Game
//	@Failure	400		{object}	ErrorResponse
//	@Failure	404		{object}	ErrorResponse
//	@Failure	409		{object}	ErrorResponse
//	@Router		/api/games/{gameID}/join [post]
func (h *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := userIDFromContext(ctx)
	gameID := chi.URLParam(r, "gameID")

	h.logger.Debug("request", "method", r.Method, "path", r.URL.Path, "gameID", gameID, "userID", userID)

	g, err := h.svc.JoinGame(ctx, userID, game.JoinGameCommand{GameID: gameID})
	if err != nil {
		h.logger.Warn("join game failed", "gameID", gameID, "userID", userID, "error", err)
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, g)
}

// JoinAnyGame godoc
//
//	@Summary	Join any available game
//	@ID		joinAnyGame
//	@Tags		games
//	@Produce	json
//	@Success	200		{object}	game.Game
//	@Failure	400		{object}	ErrorResponse
//	@Failure	404		{object}	ErrorResponse
//	@Failure	409		{object}	ErrorResponse
//	@Router		/api/games/join [post]
func (h *GameHandler) JoinAnyGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := userIDFromContext(ctx)

	h.logger.Debug("request", "method", r.Method, "path", r.URL.Path, "userID", userID)

	g, err := h.svc.JoinAnyGame(ctx, userID)
	if err != nil {
		h.logger.Warn("join any game failed", "userID", userID, "error", err)
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, g)
}

// MakeMove godoc
//
//	@Summary	Make a move
//	@ID		makeMove
//	@Tags		games
//	@Accept		json
//	@Produce	json
//	@Param		gameID	path		string				true	"Game ID"
//	@Param		request	body		game.MakeMoveCommand	true	"Move request"
//	@Success	200		{object}	game.Game
//	@Failure	400		{object}	ErrorResponse
//	@Failure	404		{object}	ErrorResponse
//	@Failure	409		{object}	ErrorResponse
//	@Router		/api/games/{gameID}/move [post]
func (h *GameHandler) MakeMove(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := userIDFromContext(ctx)
	gameID := chi.URLParam(r, "gameID")
	var cmd game.MakeMoveCommand
	if err := parseRequestBody(r, &cmd, h.logger, w); err != nil {
		return // parseRequestBody already wrote the error response
	}
	cmd.GameID = gameID

	h.logger.Debug("request", "method", r.Method, "path", r.URL.Path, "gameID", gameID, "userID", userID, "x", cmd.X, "y", cmd.Y)

	g, err := h.svc.MakeMove(ctx, userID, cmd)
	if err != nil {
		h.logger.Warn("make move failed", "gameID", gameID, "userID", userID, "x", cmd.X, "y", cmd.Y, "error", err)
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, g)
}

// GiveUp godoc
//
//	@Summary	Give up the game
//	@ID		giveUpGame
//	@Tags		games
//	@Produce	json
//	@Param		gameID	path		string	true	"Game ID"
//	@Success	200		{object}	game.Game
//	@Failure	400		{object}	ErrorResponse
//	@Failure	404		{object}	ErrorResponse
//	@Failure	409		{object}	ErrorResponse
//	@Router		/api/games/{gameID}/giveup [post]
func (h *GameHandler) GiveUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := userIDFromContext(ctx)
	gameID := chi.URLParam(r, "gameID")

	h.logger.Debug("request", "method", r.Method, "path", r.URL.Path, "gameID", gameID, "userID", userID)

	g, err := h.svc.GiveUp(ctx, userID, game.GiveUpCommand{GameID: gameID})
	if err != nil {
		h.logger.Warn("give up failed", "gameID", gameID, "userID", userID, "error", err)
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, g)
}

// Quit godoc
//
//	@Summary	Quit a game that is still waiting for an opponent
//	@ID		quitGame
//	@Tags		games
//	@Produce	json
//	@Param		gameID	path		string	true	"Game ID"
//	@Success	200		{object}	game.Game
//	@Failure	400		{object}	ErrorResponse
//	@Failure	404		{object}	ErrorResponse
//	@Failure	409		{object}	ErrorResponse
//	@Router		/api/games/{gameID}/quit [post]
func (h *GameHandler) Quit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := userIDFromContext(ctx)
	gameID := chi.URLParam(r, "gameID")

	h.logger.Debug("request", "method", r.Method, "path", r.URL.Path, "gameID", gameID, "userID", userID)

	g, err := h.svc.Quit(ctx, userID, game.QuitCommand{GameID: gameID})
	if err != nil {
		h.logger.Warn("quit failed", "gameID", gameID, "userID", userID, "error", err)
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, g)
}

// GetGame godoc
//
//	@Summary	Get game state
//	@ID		getGame
//	@Tags		games
//	@Produce	json
//	@Param		gameID	path		string	true	"Game ID"
//	@Success	200		{object}	game.Game
//	@Failure	404		{object}	ErrorResponse
//	@Router		/api/games/{gameID} [get]
func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := userIDFromContext(ctx)
	gameID := chi.URLParam(r, "gameID")

	g, err := h.svc.GetGame(ctx, userID, game.GetGameCommand{GameID: gameID})
	if err != nil {
		h.logger.Warn("get game failed", "gameID", gameID, "error", err)
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, g)
}

// GetLatestGame godoc
//
//	@Summary	Get the authenticated user's most recent game
//	@ID		getLatestGame
//	@Tags		games
//	@Produce	json
//	@Success	200	{object}	game.Game
//	@Failure	404	{object}	ErrorResponse
//	@Router		/api/games [get]
func (h *GameHandler) GetLatestGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := userIDFromContext(ctx)

	h.logger.Debug("request", "method", r.Method, "path", r.URL.Path, "userID", userID)

	g, err := h.svc.GetLatestGameForUser(ctx, userID)
	if err != nil {
		h.logger.Warn("get latest game failed", "userID", userID, "error", err)
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, g)
}
