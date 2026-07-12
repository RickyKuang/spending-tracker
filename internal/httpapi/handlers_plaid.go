package httpapi

import "net/http"

// LinkTokenCreate starts a new Plaid Link session for the caller, returning {link_session_id, link_url}.
func (h *Handlers) LinkTokenCreate(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

// LinkExchange exchanges a Link public token for a permanent access token and stores the institution.
func (h *Handlers) LinkExchange(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

// LinkStatus reports the status of a pending Link session, returning {status, institution?}.
func (h *Handlers) LinkStatus(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

// ServeLinkPage serves web/plaid_link.html, the browser-driven Plaid Link page.
func (h *Handlers) ServeLinkPage(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}
