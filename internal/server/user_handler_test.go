package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/funduck/tic-tac-toe/internal/auth"
	"github.com/funduck/tic-tac-toe/internal/user"
)

type mockUserService struct {
	signupFunc       func(ctx context.Context, userID, password string) (*user.User, *auth.TokenPair, error)
	loginFunc        func(ctx context.Context, userID, password string) (*user.User, *auth.TokenPair, error)
	refreshTokenFunc func(ctx context.Context, refreshToken string) (*user.User, *auth.TokenPair, error)
}

func (m *mockUserService) Signup(ctx context.Context, userID, password string) (*user.User, *auth.TokenPair, error) {
	return m.signupFunc(ctx, userID, password)
}

func (m *mockUserService) Login(ctx context.Context, userID, password string) (*user.User, *auth.TokenPair, error) {
	return m.loginFunc(ctx, userID, password)
}

func (m *mockUserService) RefreshToken(ctx context.Context, refreshToken string) (*user.User, *auth.TokenPair, error) {
	return m.refreshTokenFunc(ctx, refreshToken)
}

// checkAccessTokenResponse asserts the access token is in the body and that the
// refresh token is delivered only as an HttpOnly cookie, never in the body.
func checkAccessTokenResponse(wantAccess, wantRefresh string) func(t *testing.T, w *httptest.ResponseRecorder) {
	return func(t *testing.T, w *httptest.ResponseRecorder) {
		var resp AccessTokenResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if resp.AccessToken != wantAccess {
			t.Errorf("unexpected access token: %q", resp.AccessToken)
		}
		// The refresh token must not leak into the JSON body.
		var raw map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if _, ok := raw["refresh_token"]; ok {
			t.Error("refresh_token must not be present in the response body")
		}

		c := findCookie(w, refreshTokenCookieName)
		if c == nil {
			t.Fatal("expected refresh_token cookie to be set")
		}
		if c.Value != wantRefresh {
			t.Errorf("unexpected refresh cookie value: %q", c.Value)
		}
		if !c.HttpOnly {
			t.Error("refresh_token cookie must be HttpOnly")
		}
	}
}

func TestUserHandler_Signup(t *testing.T) {
	for _, tt := range []struct {
		name           string
		requestBody    UserSignupRequest
		signupFunc     func(ctx context.Context, userID, password string) (*user.User, *auth.TokenPair, error)
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "successful signup",
			requestBody: UserSignupRequest{
				UserID:   "testuser",
				Password: "testpass",
			},
			signupFunc: func(ctx context.Context, userID, password string) (*user.User, *auth.TokenPair, error) {
				return &user.User{ID: userID}, &auth.TokenPair{AccessToken: "access", RefreshToken: "refresh"}, nil
			},
			expectedStatus: http.StatusOK,
			checkResponse:  checkAccessTokenResponse("access", "refresh"),
		},
		{
			name: "signup error",
			requestBody: UserSignupRequest{
				UserID:   "testuser",
				Password: "testpass",
			},
			signupFunc: func(ctx context.Context, userID, password string) (*user.User, *auth.TokenPair, error) {
				return nil, nil, user.ErrUserAlreadyExists
			},
			expectedStatus: http.StatusConflict,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			h := NewUserHandler(&mockUserService{
				signupFunc: tt.signupFunc,
			}, slog.Default(), false)

			req := makeJSONRequest(http.MethodPost, "/api/users/signup", tt.requestBody)
			w := executeRequest(h.Signup, req)

			checkStatusCode(tt.expectedStatus)(t, w)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	for _, tt := range []struct {
		name           string
		requestBody    UserLoginRequest
		loginFunc      func(ctx context.Context, userID, password string) (*user.User, *auth.TokenPair, error)
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "successful login",
			requestBody: UserLoginRequest{
				UserID:   "testuser",
				Password: "testpass",
			},
			loginFunc: func(ctx context.Context, userID, password string) (*user.User, *auth.TokenPair, error) {
				return &user.User{ID: userID}, &auth.TokenPair{AccessToken: "access", RefreshToken: "refresh"}, nil
			},
			expectedStatus: http.StatusOK,
			checkResponse:  checkAccessTokenResponse("access", "refresh"),
		},
		{
			name: "login error",
			requestBody: UserLoginRequest{
				UserID:   "testuser",
				Password: "testpass",
			},
			loginFunc: func(ctx context.Context, userID, password string) (*user.User, *auth.TokenPair, error) {
				return nil, nil, user.ErrInvalidCredentials
			},
			expectedStatus: http.StatusUnauthorized,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			h := NewUserHandler(&mockUserService{
				loginFunc: tt.loginFunc,
			}, slog.Default(), false)

			req := makeJSONRequest(http.MethodPost, "/api/users/login", tt.requestBody)
			w := executeRequest(h.Login, req)

			checkStatusCode(tt.expectedStatus)(t, w)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestUserHandler_RefreshToken(t *testing.T) {
	for _, tt := range []struct {
		name             string
		cookie           *http.Cookie
		refreshTokenFunc func(ctx context.Context, refreshToken string) (*user.User, *auth.TokenPair, error)
		expectedStatus   int
		checkResponse    func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:   "successful token refresh",
			cookie: &http.Cookie{Name: refreshTokenCookieName, Value: "valid-refresh-token"},
			refreshTokenFunc: func(ctx context.Context, refreshToken string) (*user.User, *auth.TokenPair, error) {
				return &user.User{ID: "testuser"}, &auth.TokenPair{AccessToken: "new-access", RefreshToken: "new-refresh"}, nil
			},
			expectedStatus: http.StatusOK,
			checkResponse:  checkAccessTokenResponse("new-access", "new-refresh"),
		},
		{
			name:   "invalid refresh token",
			cookie: &http.Cookie{Name: refreshTokenCookieName, Value: "invalid-refresh-token"},
			refreshTokenFunc: func(ctx context.Context, refreshToken string) (*user.User, *auth.TokenPair, error) {
				return nil, nil, auth.ErrTokenInvalid
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "missing refresh cookie",
			cookie:         nil,
			expectedStatus: http.StatusUnauthorized,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			h := NewUserHandler(&mockUserService{
				refreshTokenFunc: tt.refreshTokenFunc,
			}, slog.Default(), false)

			req := makeJSONRequest(http.MethodPost, "/api/users/refresh-token", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}
			w := executeRequest(h.RefreshToken, req)

			checkStatusCode(tt.expectedStatus)(t, w)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}
