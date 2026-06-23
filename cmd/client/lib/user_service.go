package lib

import (
	"context"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID"
)

// UserService wraps the generated API client with convenience methods
type UserService struct {
	api *openapi.UsersAPIService
}

// NewUserService creates a new UserService
func NewUserService(apiClient *openapi.APIClient) *UserService {
	return &UserService{api: apiClient.UsersAPI}
}

// Signup registers a new user and returns the access token. The refresh token is
// returned as an HttpOnly cookie and stored by the HTTP client's cookie jar.
func (s *UserService) Signup(ctx context.Context, userID, password string) (*openapi.ServerAccessTokenResponse, error) {
	req := openapi.NewServerUserSignupRequest()
	req.SetUserId(userID)
	req.SetPassword(password)
	resp, r, err := s.api.Signup(ctx).Request(*req).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return resp, nil
}

// Login authenticates an existing user and returns the access token. The refresh
// token is returned as an HttpOnly cookie and stored by the cookie jar.
func (s *UserService) Login(ctx context.Context, userID, password string) (*openapi.ServerAccessTokenResponse, error) {
	req := openapi.NewServerUserLoginRequest()
	req.SetUserId(userID)
	req.SetPassword(password)
	resp, r, err := s.api.Login(ctx).Request(*req).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return resp, nil
}

// RefreshToken exchanges the refresh token (sent automatically as an HttpOnly
// cookie by the cookie jar) for a new access token.
func (s *UserService) RefreshToken(ctx context.Context) (*openapi.ServerAccessTokenResponse, error) {
	resp, r, err := s.api.RefreshToken(ctx).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return resp, nil
}
