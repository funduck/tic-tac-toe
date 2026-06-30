package lib

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
)

// AuthRoundTripper is an http.RoundTripper that adds authentication headers
// and handles token refresh on 401 errors
type AuthRoundTripper struct {
	transport   http.RoundTripper
	userService *UserService
	mu          sync.Mutex // Protects concurrent token refresh operations
}

// NewAuthRoundTripper creates a new AuthRoundTripper
func NewAuthRoundTripper(userService *UserService) *AuthRoundTripper {
	return &AuthRoundTripper{
		transport:   http.DefaultTransport,
		userService: userService,
	}
}

// RoundTrip executes a single HTTP transaction with authentication
func (a *AuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	reqClone := cloneRequest(req)

	// The refresh-token request must never recurse through the refresh-on-401
	// logic below: it would re-enter RoundTrip and deadlock on a.mu (already held
	// by the in-flight refresh). Send it straight through.
	if strings.HasSuffix(req.URL.Path, "/api/users/refresh-token") {
		return a.transport.RoundTrip(reqClone)
	}

	// Add Authorization header if we have an access token
	if AccessToken != "" {
		reqClone.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AccessToken))
	}

	// Execute the request
	resp, err := a.transport.RoundTrip(reqClone)
	if err != nil {
		return resp, err
	}

	// Check if we got a 401 Unauthorized response. We only attempt a refresh if we
	// were authenticated; the refresh token itself rides along as an HttpOnly cookie.
	if resp.StatusCode == http.StatusUnauthorized && AccessToken != "" {
		// Close the original response body
		resp.Body.Close()

		// Try to refresh the token
		if err := a.refreshToken(req.Context()); err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}

		// Retry the original request with the new token
		reqRetry := cloneRequest(req)
		reqRetry.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AccessToken))
		return a.transport.RoundTrip(reqRetry)
	}

	return resp, nil
}

// refreshToken attempts to refresh the access token using the refresh token
func (a *AuthRoundTripper) refreshToken(ctx context.Context) error {
	// Use mutex to prevent concurrent token refresh attempts
	a.mu.Lock()
	defer a.mu.Unlock()

	// The refresh token is sent automatically as an HttpOnly cookie by the cookie jar.
	resp, err := a.userService.RefreshToken(ctx)
	if err != nil {
		return err
	}

	// Update the access token; the new refresh token is stored in the cookie jar.
	AccessToken = *resp.AccessToken

	return nil
}

// cloneRequest creates a shallow copy of the request with a cloned body
func cloneRequest(req *http.Request) *http.Request {
	reqClone := req.Clone(req.Context())

	// If the request has a body, we need to clone it
	if req.Body != nil {
		// Read the body
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			// If we can't read the body, just use the original request
			return req
		}
		// Restore the original body
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		// Set the cloned body
		reqClone.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return reqClone
}

// NewAuthHTTPClient creates an HTTP client with authentication support. The cookie
// jar transparently stores the HttpOnly refresh token cookie returned on login and
// sends it back on refresh requests, so the refresh token never lives in process state.
func NewAuthHTTPClient(userService *UserService) *http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		// cookiejar.New(nil) never returns an error, but fail safe just in case.
		panic(fmt.Sprintf("failed to create cookie jar: %v", err))
	}
	return &http.Client{
		Transport: NewAuthRoundTripper(userService),
		Jar:       jar,
	}
}
