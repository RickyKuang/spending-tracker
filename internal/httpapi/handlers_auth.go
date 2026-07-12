package httpapi

import "net/http"

// SignUp creates a new user and an initial session, returning {user, session_token}.
func (h *Handlers) SignUp(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

// SignIn validates credentials and creates a new session, returning {user, session_token}.
func (h *Handlers) SignIn(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

// SignOut revokes the caller's current session.
func (h *Handlers) SignOut(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}
