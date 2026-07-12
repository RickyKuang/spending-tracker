// Command api is the spending-tracker backend HTTP server.
package main

import (
	"log/slog"
	"os"

	"github.com/RickyKuang/spending-tracker/internal/config"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	_ = cfg

	// TODO: implement
	//
	// This is the composition root: everything gets built here and nowhere else.
	//   1. Run goose migrations against cfg.DatabaseURL (see internal/store.NewPool
	//      for the pgxpool, but goose needs a database/sql.DB — open one separately
	//      with the pgx stdlib driver).
	//   2. Open a pgxpool.Pool via store.NewPool, wrap it with store.New, and call
	//      ResetStuckSyncs (crash-recovery: clears any institution left mid-sync).
	//   3. Build the leaf dependencies: crypto.NewSealer(cfg.AppEncryptionKey),
	//      plaid.New(cfg.PlaidClientID, cfg.PlaidSecret, cfg.PlaidEnv).
	//   4. Build the services that depend on those: sync.New(store, plaidClient,
	//      sealer), auth.NewSessionManager(store, ttl).
	//   5. Build httpapi.Handlers from all of the above and pass it to
	//      httpapi.NewRouter to get the http.Handler.
	//   6. Construct an *http.Server{Addr: cfg.HTTPAddr, Handler: router} and run
	//      ListenAndServe in a goroutine (log+os.Exit on unexpected errors, but
	//      ignore http.ErrServerClosed — that's the expected shutdown signal).
	//   7. Block on signal.NotifyContext(context.Background(), os.Interrupt,
	//      syscall.SIGTERM) until a signal arrives, then call srv.Shutdown with a
	//      bounded context.WithTimeout so in-flight requests can finish.
}
