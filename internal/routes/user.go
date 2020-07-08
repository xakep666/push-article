package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"push_article/internal/middlewares"

	"github.com/go-chi/chi"
)

type UserService struct{}

func (u *UserService) login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request struct {
		ID uint64 `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "example-push-user",
		Value:    strconv.FormatUint(request.ID, 10),
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour),
		Path:     "/",
	})
}

func (u *UserService) logout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	http.SetCookie(w, &http.Cookie{
		Name:     "example-push-user",
		HttpOnly: true,
		Expires:  time.Now(),
		Path:     "/",
	})
}

func (u *UserService) me(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": middlewares.UserFromContext(r.Context()),
	})
}

func (u *UserService) AddToRouter(r chi.Router) {
	r.Post("/login", u.login)
	r.With(middlewares.Authorized).Post("/logout", u.logout)
	r.With(middlewares.Authorized).Get("/me", u.me)
}
