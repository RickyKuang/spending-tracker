package httpapi

import "net/http"

// ListInstitutions lists every institution linked by the caller.
func (h *Handlers) ListInstitutions(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

// ListAccounts lists every account across the caller's linked institutions.
func (h *Handlers) ListAccounts(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}
