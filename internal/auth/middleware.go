package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/lordaris/gotth-boilerplate/internal/app"
	"github.com/lordaris/gotth-boilerplate/internal/models"
)

type contextKey string

const UserContextKey = contextKey("user")

type Middleware struct {
	App *app.Application
}

func NewMiddleware(app *app.Application) *Middleware {
	return &Middleware{App: app}
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := m.App.DB
		// Add Vary header to ensure cached responses are based on the Authorization header
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		token := headerParts[1]

		if err := ValidateTokenPlaintext(token); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		tokenHash := HashToken(token)

		var userToken models.Token
		var user models.User

		query := `
            SELECT id, user_id, expiry 
            FROM tokens 
            WHERE token_hash = $1 
            AND expiry > $2
        `
		err := db.Get(&userToken, query, tokenHash, time.Now())
		if err != nil {
			http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		query = `
            SELECT id, username, password_hash, created_at
            FROM users
            WHERE id = $1
        `
		err = db.Get(&user, query, userToken.UserID)
		if err != nil {
			http.Error(w, "Invalid user", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, &user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r.Context())
		if user == nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserFromContext(ctx context.Context) *models.User {
	user, ok := ctx.Value(UserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}

func (m *Middleware) CleanupExpiredTokens() error {
	db := m.App.DB
	query := `DELETE FROM tokens WHERE expiry < $1`
	_, err := db.Exec(query, time.Now())
	return err
}
