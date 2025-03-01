package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/lordaris/gotth-boilerplate/internal/auth"
	"github.com/lordaris/gotth-boilerplate/internal/models"

	"github.com/lordaris/gotth-boilerplate/internal/templates"
)

func (h *Handlers) LoginForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserFromContext(r.Context())
		if user != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		templates.LoginPage().Render(r.Context(), w)
	}
}

func (h *Handlers) RegisterForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserFromContext(r.Context())
		if user != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		templates.RegisterPage().Render(r.Context(), w)
	}
}

func (h *Handlers) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := h.App.DB

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			templates.LoginPage().Render(r.Context(), w)
			return
		}

		var user models.User
		query := `SELECT id, username, password_hash, created_at FROM users WHERE username = $1`
		err := db.Get(&user, query, username)
		if err != nil {
			templates.LoginPage().Render(r.Context(), w)
			return
		}

		match, err := user.CheckPassword(password)
		if err != nil || !match {
			templates.LoginPage().Render(r.Context(), w)
			return
		}

		tokenString, tokenHash, expiryTime, err := auth.GenerateToken(24 * time.Hour)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		query = `
            INSERT INTO tokens (user_id, token_hash, plaintext_token, expiry) 
            VALUES ($1, $2, $3, $4)
        `
		_, err = db.Exec(query, user.ID, tokenHash, tokenString, expiryTime)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    tokenString,
			Path:     "/",
			Expires:  expiryTime,
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteStrictMode,
		})


		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/")
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *Handlers) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := h.App.DB

		if err := r.ParseForm(); err != nil {
			log.Println("Error while parsing the form:", err)
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		if username == "" || password == "" || password != confirmPassword {
			log.Println("Failed validation")
			templates.RegisterPage().Render(r.Context(), w)
			return
		}

		user := models.User{
			Username: username,
		}

		if err := user.SetPassword(password); err != nil {
      log.Println("Error while setting the password:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if db == nil {
			log.Println("Error: Database not initialized")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		query := `
            INSERT INTO users (username, password_hash) 
            VALUES ($1, $2) 
            RETURNING id, created_at
        `

		row := db.QueryRow(query, user.Username, user.PasswordHash)
		if err := row.Scan(&user.ID, &user.CreatedAt); err != nil {
			log.Println("Error while executing the query:", err)
			http.Error(w, "Failed to create the user", http.StatusInternalServerError)
			return
		}

		log.Println("User successfully registered", user.Username)

		// If HTMX request
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/login")
			return
		}

		// Regular form submission
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}


func (h *Handlers) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := h.App.DB
		
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			tokenHash := auth.HashToken(cookie.Value)
			query := `DELETE FROM tokens WHERE token_hash = $1`
			_, err := db.Exec(query, tokenHash)
			if err != nil {
				log.Printf("Error deleting token from DB: %v", err)
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "auth_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1, // Expire immediately
				HttpOnly: true,
				Secure:   r.TLS != nil,
				SameSite: http.SameSiteStrictMode,
			})
		}

		// If HTMX request, set redirect header
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/") 	 // Redirect to main
      return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

