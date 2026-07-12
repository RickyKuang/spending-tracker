package httpapi

import "net/http"

// ListTransactions lists the caller's transactions, filtered by account/date range and paginated by cursor.
func (h *Handlers) ListTransactions(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}
