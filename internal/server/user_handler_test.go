package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"testing"

	"github.com/funduck/tic-tac-toe/internal/auth"
	"github.com/funduck/tic-tac-toe/internal/user"
)

type mockUserService struct {
	signupFunc       func(userID, password string) (*user.User, *auth.TokenPair, error)
	loginFunc        func(userID, password string) (*user.User, *auth.TokenPair, error)
	refreshTokenFunc func(userID, refreshToken string) (*user.User, *auth.TokenPair, error)
}

func (m *mockUserService) Signup(userID, password string) (*user.User, *auth.TokenPair, error) {
	return m.signupFunc(userID, password)
}

func (m *mockUserService) Login(userID, password string) (*user.User, *auth.TokenPair, error) {
	return m.loginFunc(userID, password)
}

func (m *mockUserService) RefreshToken(userID, refreshToken string) (*user.User, *auth.TokenPair, error) {
	return m.refreshTokenFunc(userID, refreshToken)
}

func TestUserHandler_Signup(t *testing.T) {
	for _, tt := range []struct {
		name           string
		requestBody    UserSignupRequest
		signupFunc     func(userID, password string) (*user.User, *auth.TokenPair, error)
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successful signup",
			requestBody: UserSignupRequest{
				UserID:   "testuser",
				Password: "testpass",
			},
			signupFunc: func(userID, password string) (*user.User, *auth.TokenPair, error) {
				return &user.User{ID: userID}, &auth.TokenPair{AccessToken: "access", RefreshToken: "refresh"}, nil
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var resp auth.TokenPair
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.AccessToken != "access" || resp.RefreshToken != "refresh" {
					t.Errorf("unexpected tokens: %+v", resp)
				}
			},
		},
		{
			name: "signup error",
			requestBody: UserSignupRequest{
				UserID:   "testuser",
				Password: "testpass",
			},
			signupFunc: func(userID, password string) (*user.User, *auth.TokenPair, error) {
				return nil, nil, user.ErrUserAlreadyExists
			},
			expectedStatus: http.StatusConflict,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			h := NewUserHandler(&mockUserService{
				signupFunc: tt.signupFunc,
			}, slog.Default())

			req := makeJSONRequest(http.MethodPost, "/api/users/signup", tt.requestBody)
			w := executeRequest(h.Signup, req)

			checkStatusCode(tt.expectedStatus)(t, w)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	for _, tt := range []struct {
		name           string
		requestBody    UserLoginRequest
		loginFunc      func(userID, password string) (*user.User, *auth.TokenPair, error)
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successful login",
			requestBody: UserLoginRequest{
				UserID:   "testuser",
				Password: "testpass",
			},
			loginFunc: func(userID, password string) (*user.User, *auth.TokenPair, error) {
				return &user.User{ID: userID}, &auth.TokenPair{AccessToken: "access", RefreshToken: "refresh"}, nil
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var resp auth.TokenPair
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.AccessToken != "access" || resp.RefreshToken != "refresh" {
					t.Errorf("unexpected tokens: %+v", resp)
				}
			},
		},
		{
			name: "login error",
			requestBody: UserLoginRequest{
				UserID:   "testuser",
				Password: "testpass",
			},
			loginFunc: func(userID, password string) (*user.User, *auth.TokenPair, error) {
				return nil, nil, user.ErrInvalidCredentials
			},
			expectedStatus: http.StatusUnauthorized,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			h := NewUserHandler(&mockUserService{
				loginFunc: tt.loginFunc,
			}, slog.Default())

			req := makeJSONRequest(http.MethodPost, "/api/users/login", tt.requestBody)
			w := executeRequest(h.Login, req)

			checkStatusCode(tt.expectedStatus)(t, w)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestUserHandler_RefreshToken(t *testing.T) {
	for _, tt := range []struct {
		name             string
		requestBody      UserRefreshTokenRequest
		refreshTokenFunc func(userID, refreshToken string) (*user.User, *auth.TokenPair, error)
		expectedStatus   int
		checkResponse    func(t *testing.T, body []byte)
	}{
		{
			name: "successful token refresh",
			requestBody: UserRefreshTokenRequest{
				UserID:       "testuser",
				RefreshToken: "valid-refresh-token",
			},
			refreshTokenFunc: func(userID, refreshToken string) (*user.User, *auth.TokenPair, error) {
				return &user.User{ID: userID}, &auth.TokenPair{AccessToken: "new-access", RefreshToken: "new-refresh"}, nil
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var resp auth.TokenPair
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.AccessToken != "new-access" || resp.RefreshToken != "new-refresh" {
					t.Errorf("unexpected tokens: %+v", resp)
				}
			},
		},
		{
			name: "invalid refresh token",
			requestBody: UserRefreshTokenRequest{
				UserID:       "testuser",
				RefreshToken: "invalid-refresh-token",
			},
			refreshTokenFunc: func(userID, refreshToken string) (*user.User, *auth.TokenPair, error) {
				return nil, nil, auth.ErrTokenInvalid
			},
			expectedStatus: http.StatusUnauthorized,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			h := NewUserHandler(&mockUserService{
				refreshTokenFunc: tt.refreshTokenFunc,
			}, slog.Default())

			req := makeJSONRequest(http.MethodPost, "/api/users/refresh-token", tt.requestBody)
			w := executeRequest(h.RefreshToken, req)

			checkStatusCode(tt.expectedStatus)(t, w)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}
