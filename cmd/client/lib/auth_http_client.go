package lib

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
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

	// Add Authorization header if we have an access token
	if AccessToken != "" {
		reqClone.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AccessToken))
	}

	// Execute the request
	resp, err := a.transport.RoundTrip(reqClone)
	if err != nil {
		return resp, err
	}

	// Check if we got a 401 Unauthorized response
	if resp.StatusCode == http.StatusUnauthorized && RefreshToken != "" {
		// Close the original response body
		resp.Body.Close()

		log.Printf("Received 401 Unauthorized. Attempting to refresh token...")

		// Try to refresh the token
		if err := a.refreshToken(req.Context()); err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}

		log.Printf("Token refreshed successfully. Retrying original request...")

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

	// Double-check if another goroutine already refreshed the token
	// by attempting the refresh with current refresh token
	tokenPair, err := a.userService.RefreshToken(ctx, UserID, RefreshToken)
	if err != nil {
		return fmt.Errorf("refresh token request failed: %w", err)
	}

	// Update the global tokens
	AccessToken = *tokenPair.AccessToken
	RefreshToken = *tokenPair.RefreshToken

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

// NewAuthHTTPClient creates an HTTP client with authentication support
func NewAuthHTTPClient(userService *UserService) *http.Client {
	return &http.Client{
		Transport: NewAuthRoundTripper(userService),
	}
}
