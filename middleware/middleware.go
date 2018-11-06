package middleware

import (
	"net/http"
	"../sessions"
)

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// before action
		session, _ := sessions.Store.Get(r, "sessionId")
		_, ok := session.Values["user"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		handler.ServeHTTP(w, r)
	}
}