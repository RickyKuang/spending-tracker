package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/RickyKuang/spending-tracker/internal/auth"
	"github.com/RickyKuang/spending-tracker/internal/bank"
	"github.com/RickyKuang/spending-tracker/internal/store"
	"github.com/RickyKuang/spending-tracker/internal/sync"
)

// Handlers holds the dependencies shared by every HTTP handler.
type Handlers struct {
	Store    *store.PostgresStore
	SyncSvc  *sync.Service
	Sessions *auth.SessionManager
	Bank     bank.Provider
	Logger   *slog.Logger
}

// NewRouter builds the chi router with every route wired to its handler, per the endpoint
// inventory in docs/ARCHITECTURE.md.
func NewRouter(h *Handlers) http.Handler {
	r := chi.NewRouter()

	r.Get("/healthz", healthz)
	r.Get("/plaid/link", h.ServeLinkPage)

	r.Post("/api/auth/signup", h.SignUp)
	r.Post("/api/auth/signin", h.SignIn)
	r.Post("/api/plaid/link/exchange", h.LinkExchange)

	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(h.Sessions))

		r.Post("/api/auth/signout", h.SignOut)
		r.Post("/api/plaid/link/token", h.LinkTokenCreate)
		r.Get("/api/plaid/link/status/{id}", h.LinkStatus)
		r.Get("/api/institutions", h.ListInstitutions)
		r.Get("/api/accounts", h.ListAccounts)
		r.Get("/api/transactions", h.ListTransactions)
		r.Post("/api/sync", h.Sync)
	})

	return r
}

// healthz reports process liveness for operational checks.
func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
