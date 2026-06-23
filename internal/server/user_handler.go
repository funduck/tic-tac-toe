package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/funduck/tic-tac-toe/internal/auth"
	"github.com/funduck/tic-tac-toe/internal/user"
)

type UserService interface {
	Signup(ctx context.Context, userID, password string) (*user.User, *auth.TokenPair, error)
	Login(ctx context.Context, userID, password string) (*user.User, *auth.TokenPair, error)
	RefreshToken(ctx context.Context, userID, refreshToken string) (*user.User, *auth.TokenPair, error)
}

type UserHandler struct {
	svc    UserService
	logger *slog.Logger
}

func NewUserHandler(svc UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		svc:    svc,
		logger: logger,
	}
}

// Signup godoc
//
//	@Summary	Sign up a new user
//	@ID		signup
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		request	body		UserSignupRequest	true	"User signup request"
//	@Success	200		{object}	auth.TokenPair
//	@Failure	400		{object}	ErrorResponse
//	@Failure	409		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/api/users/signup [post]
func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req UserSignupRequest
	if err := parseRequestBody(r, &req, h.logger, w); err != nil {
		return
	}

	_, tokens, err := h.svc.Signup(ctx, req.UserID, req.Password)
	if err != nil {
		h.logger.Warn("signup failed", "userID", req.UserID, "error", err)
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

// Login godoc
//
//	@Summary	Log in an existing user
//	@ID		login
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		request	body		UserLoginRequest	true	"User login request"
//	@Success	200		{object}	auth.TokenPair
//	@Failure	400		{object}	ErrorResponse
//	@Failure	401		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/api/users/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req UserLoginRequest
	if err := parseRequestBody(r, &req, h.logger, w); err != nil {
		return
	}

	_, tokens, err := h.svc.Login(ctx, req.UserID, req.Password)
	if err != nil {
		h.logger.Warn("login failed", "userID", req.UserID, "error", err)
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

// RefreshToken godoc
//
//	@Summary	Refresh access token using refresh token
//	@ID		refresh-token
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		request	body		UserRefreshTokenRequest	true	"User refresh token request"
//	@Success	200		{object}	auth.TokenPair
//	@Failure	400		{object}	ErrorResponse
//	@Failure	401		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/api/users/refresh-token [post]
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req UserRefreshTokenRequest
	if err := parseRequestBody(r, &req, h.logger, w); err != nil {
		return
	}

	_, tokens, err := h.svc.RefreshToken(ctx, req.UserID, req.RefreshToken)
	if err != nil {
		h.logger.Warn("refresh token failed", "userID", req.UserID, "error", err)
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}
