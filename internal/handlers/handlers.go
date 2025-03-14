package handlers

import (
	"net/http"

	// "github.com/go-chi/chi/v5"
	"github.com/lordaris/gotth-boilerplate/internal/app"
	"github.com/lordaris/gotth-boilerplate/internal/auth"
	"github.com/lordaris/gotth-boilerplate/internal/models"
	"github.com/lordaris/gotth-boilerplate/internal/templates"
)

type Handlers struct {
	App *app.Application
}

func NewHandlers(app *app.Application) *Handlers {
	return &Handlers{
		App: app,
	}
}

func (h *Handlers) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserFromContext(r.Context())
		if user == nil {
			user = &models.User{}
		}
		templates.Home(*user).Render(r.Context(), w)
	}
}
//
func (h *Handlers) ProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserFromContext(r.Context())
		if user == nil {
			user = &models.User{}
		}
		templates.ProfilePage(user).Render(r.Context(), w)
	}
}

// func (h *Handlers) CreateUserHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		users := []models.User{}
// 		err := r.ParseForm()
// 		if err != nil {
// 			http.Error(w, "Invalid input", http.StatusBadRequest)
// 			return
// 		}
//
// 		name := r.FormValue("name")
// 		email := r.FormValue("email")
//
// 		if name == "" || email == "" {
// 			http.Error(w, "Missing fields", http.StatusBadRequest)
// 			return
// 		}
//
// 		_, err = h.App.DB.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
// 		if err != nil {
// 			http.Error(w, "Error creating user", http.StatusInternalServerError)
// 			return
// 		}
//
// 		err = h.App.DB.Select(&users, "SELECT id, name, email FROM users")
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		w.WriteHeader(http.StatusCreated)
// 		templates.UsersList(users).Render(r.Context(), w)
// 	}
// }
//
// // GetUsersHandler retrieves users from the database
// func (h *Handlers) GetUsersHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		users := []models.User{}
// 		err := h.App.DB.Select(&users, "SELECT id, name, email FROM users")
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		templates.UsersList(users).Render(r.Context(), w)
// 	}
// }
//
// func (h *Handlers) GetUserDetailHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		idStr := chi.URLParam(r, "id")
// 		id, err := strconv.Atoi(idStr)
// 		if err != nil {
// 			log.Println("Error converting ID:", err)
// 			http.Error(w, "Invalid user ID", http.StatusBadRequest)
// 			return
// 		}
//
// 		user := models.User{}
// 		err = h.App.DB.Get(&user, "SELECT id, name, email FROM users WHERE id = $1", id)
// 		if err != nil {
// 			log.Println("Error fetching user from DB:", err)
// 			http.Error(w, "User not found", http.StatusNotFound)
// 			return
// 		}
//
// 		templates.UserDetails(user).Render(r.Context(), w)
// 	}
// }
//
// func (h *Handlers) EditUserHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		idStr := chi.URLParam(r, "id")
// 		id, err := strconv.Atoi(idStr)
// 		if err != nil {
// 			http.Error(w, "Invalid user ID", http.StatusBadRequest)
// 			return
// 		}
//
// 		user := models.User{}
// 		err = h.App.DB.Get(&user, "SELECT id, name, email FROM users WHERE id = $1", id)
// 		if err != nil {
// 			http.Error(w, "User not found", http.StatusNotFound)
// 			return
// 		}
//
// 		templates.EditUserForm(user).Render(r.Context(), w)
// 	}
// }
//
// func (h *Handlers) UpdateUserHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		id := chi.URLParam(r, "id")
// 		err := r.ParseForm()
// 		if err != nil {
// 			http.Error(w, "Invalid input", http.StatusBadRequest)
// 			return
// 		}
//
// 		name := r.FormValue("name")
// 		email := r.FormValue("email")
//
// 		_, err = h.App.DB.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", name, email, id)
// 		if err != nil {
// 			http.Error(w, "Error updating user", http.StatusInternalServerError)
// 			return
// 		}
// 		users := []models.User{}
// 		err = h.App.DB.Select(&users, "SELECT id, name, email FROM users")
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		templates.UsersList(users).Render(r.Context(), w)
// 	}
// }
//
// func (h *Handlers) DeleteUserHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		id := chi.URLParam(r, "id")
//
// 		_, err := h.App.DB.Exec("DELETE FROM users WHERE id = $1", id)
// 		if err != nil {
// 			http.Error(w, "Error deleting user", http.StatusInternalServerError)
// 			return
// 		}
//
// 		w.WriteHeader(http.StatusOK)
//
// 		users := []models.User{}
// 		err = h.App.DB.Select(&users, "SELECT id, name, email FROM users")
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		templates.UsersList(users).Render(r.Context(), w)
// 	}
// }
