package lib

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

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

func parseUserError(r *http.Response, err error) error {
	if r == nil {
		return err
	}
	e := &openapi.ServerErrorResponse{}
	decoder := json.NewDecoder(r.Body)
	if err2 := decoder.Decode(e); err2 != nil {
		return err
	}
	return errors.Join(errors.New(e.GetError()), err)
}

func (s *UserService) Signup(ctx context.Context, userID, password string) (*openapi.AuthTokenPair, error) {
	req := openapi.NewServerUserSignupRequest()
	req.SetUserID(userID)
	req.SetPassword(password)
	tokenPair, r, err := s.api.Signup(ctx).Request(*req).Execute()
	if err != nil {
		return nil, parseUserError(r, err)
	}
	return tokenPair, nil
}

func (s *UserService) Login(ctx context.Context, userID, password string) (*openapi.AuthTokenPair, error) {
	req := openapi.NewServerUserLoginRequest()
	req.SetUserID(userID)
	req.SetPassword(password)
	tokenPair, r, err := s.api.Login(ctx).Request(*req).Execute()
	if err != nil {
		return nil, parseUserError(r, err)
	}
	return tokenPair, nil
}

func (s *UserService) RefreshToken(ctx context.Context, userID, refreshToken string) (*openapi.AuthTokenPair, error) {
	req := openapi.NewServerUserRefreshTokenRequest()
	req.SetUserID(userID)
	req.SetRefreshToken(refreshToken)
	tokenPair, r, err := s.api.RefreshToken(ctx).Request(*req).Execute()
	if err != nil {
		return nil, parseUserError(r, err)
	}
	return tokenPair, nil
}
