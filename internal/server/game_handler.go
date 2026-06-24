package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/funduck/tic-tac-toe/internal/game"
)

type GameService interface {
	CreateGame(ctx context.Context) (*game.Game, error)
	CreatePrivateGame(ctx context.Context) (*game.Game, error)
	JoinGame(ctx context.Context, gameID, userID string) (*game.Game, error)
	JoinAnyGame(ctx context.Context, userID string) (*game.Game, error)
	MakeMove(ctx context.Context, gameID, userID string, x, y int) (*game.Game, error)
	GiveUp(ctx context.Context, gameID, userID string) (*game.Game, error)
	Quit(ctx context.Context, gameID, userID string) (*game.Game, error)
	GetGame(ctx context.Context, gameID, userID string) (*game.Game, error)
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
//	@Param		request	body		CreateGameRequest	true	"Create game request"
//	@Success	201		{object}	game.Game
//	@Failure	400		{object}	ErrorResponse
//	@Failure	409		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/api/games [post]
func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := userIDFromContext(ctx)
	var req CreateGameRequest
	if err := parseRequestBody(r, &req, h.logger, w); err != nil {
		return // parseRequestBody already wrote the error response
	}

	var g *game.Game
	var err error

	if req.Private {
		g, err = h.svc.CreatePrivateGame(ctx)
	} else {
		g, err = h.svc.CreateGame(ctx)
	}

	if err != nil {
		h.logger.Warn("create game failed", "error", err)
		writeDomainError(w, err)
		return
	}

	g2, err := h.svc.JoinGame(ctx, g.ID, userID)
	if err != nil {
		h.logger.Warn("join game failed", "gameID", g.ID, "userID", userID, "error", err)
		writeDomainError(w, err)
		return
	}
	g = g2

	writeJSON(w, http.StatusCreated, g)
}

// JoinGame godoc
//
//	@Summary	Join a waiting game
//	@ID		joinGame
//	@Tags		games
//	@Accept		json
//	@Produce	json
//	@Param		gameID	path		string			true	"Game ID"
//	@Param		request	body		JoinGameRequest	true	"Join game request"
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

	g, err := h.svc.JoinGame(ctx, gameID, userID)
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
//	@Accept		json
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
//	@Param		gameID	path		string		true	"Game ID"
//	@Param		request	body		MoveRequest	true	"Move request"
//	@Success	200		{object}	game.Game
//	@Failure	400		{object}	ErrorResponse
//	@Failure	404		{object}	ErrorResponse
//	@Failure	409		{object}	ErrorResponse
//	@Router		/api/games/{gameID}/move [post]
func (h *GameHandler) MakeMove(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := userIDFromContext(ctx)
	gameID := chi.URLParam(r, "gameID")
	var req MoveRequest
	if err := parseRequestBody(r, &req, h.logger, w); err != nil {
		return // parseRequestBody already wrote the error response
	}

	h.logger.Debug("request", "method", r.Method, "path", r.URL.Path, "gameID", gameID, "userID", userID, "x", req.X, "y", req.Y)

	g, err := h.svc.MakeMove(ctx, gameID, userID, req.X, req.Y)
	if err != nil {
		h.logger.Warn("make move failed", "gameID", gameID, "userID", userID, "x", req.X, "y", req.Y, "error", err)
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
//	@Accept		json
//	@Produce	json
//	@Param		gameID	path		string			true	"Game ID"
//	@Param		request	body		GiveUpRequest	true	"Give up request"
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

	g, err := h.svc.GiveUp(ctx, gameID, userID)
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
//	@Accept		json
//	@Produce	json
//	@Param		gameID	path		string		true	"Game ID"
//	@Param		request	body		QuitRequest	true	"Quit request"
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

	g, err := h.svc.Quit(ctx, gameID, userID)
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

	g, err := h.svc.GetGame(ctx, gameID, userID)
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
