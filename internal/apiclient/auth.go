package apiclient

import (
	"context"

	"github.com/RickyKuang/spending-tracker/internal/models"
)

// SignUpResult is the response from POST /api/auth/signup.
type SignUpResult struct {
	User         models.User
	SessionToken string
}

// SignInResult is the response from POST /api/auth/signin.
type SignInResult struct {
	User         models.User
	SessionToken string
}

// SignUp creates a new user account.
func (c *Client) SignUp(ctx context.Context, email, password string) (SignUpResult, error) {
	// TODO: implement
	return SignUpResult{}, nil
}

// SignIn authenticates an existing user and returns a new session token.
func (c *Client) SignIn(ctx context.Context, email, password string) (SignInResult, error) {
	// TODO: implement
	return SignInResult{}, nil
}

// SignOut revokes the current session token.
func (c *Client) SignOut(ctx context.Context) error {
	// TODO: implement
	return nil
}
