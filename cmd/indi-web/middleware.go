package main

import (
	"net/http"

	indiserver "github.com/adriffaud/indi-web/internal/indi-server"
)

func (app *application) checkServerStarted(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !indiserver.IsRunning() && r.URL.Path != "/setup" {
			http.Redirect(w, r, "/setup", http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}
