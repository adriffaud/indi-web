package main

import (
	"log/slog"
	"net/http"
	"strings"

	indiserver "github.com/adriffaud/indi-web/internal/indi-server"
)

func (app *application) checkServerStarted(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Checking if INDI server is running")

		if !indiserver.IsRunning() && r.URL.Path != "/setup" && !strings.HasPrefix(r.URL.Path, "/static/") {
			slog.Debug("INDI server is not running, redirecting to setup")
			http.Redirect(w, r, "/setup", http.StatusTemporaryRedirect)
			return
		} else if indiserver.IsRunning() && r.URL.Path == "/setup" {
			slog.Debug("INDI server is running, redirecting to index")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}
