package routes

import (
	"net/http"

	"github.com/lordaris/gotth-boilerplate/internal/app"
	"github.com/lordaris/gotth-boilerplate/internal/auth"
	"github.com/lordaris/gotth-boilerplate/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// SetupRouter configures and returns the application router
func SetupRouter(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	authMiddleware := auth.NewMiddleware(app)

	// Use cookie-based auth by checking for auth_token cookie
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				// Add Authorization header if cookie exists
				r.Header.Set("Authorization", "Bearer "+cookie.Value)
			}
			next.ServeHTTP(w, r)
		})
	})

	// Apply authentication middleware to parse the auth token
	r.Use(authMiddleware.Authenticate)

	// Create handlers
	h := handlers.NewHandlers(app)

	// Static files
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", h.Home())
		r.Get("/login", h.LoginForm())
		r.Post("/login", h.Login())
		r.Get("/register", h.RegisterForm())
		r.Post("/register", h.Register())
		r.Get("/logout", h.Logout())
	})

	// User management routes
	r.Group(func(r chi.Router) {
		// These routes require authentication
		r.Use(authMiddleware.RequireAuth)
		r.Get("/profile", h.ProfileHandler())

		// User management
		// r.Get("/users", h.GetUsersHandler())
		// r.Post("/users", h.CreateUserHandler())
		// r.Get("/users/{id}", h.GetUserDetailHandler())
		// r.Get("/users/{id}/edit", h.EditUserHandler())
		// r.Put("/users/{id}", h.UpdateUserHandler())
		// r.Delete("/users/{id}", h.DeleteUserHandler())
	})

	return r
}
