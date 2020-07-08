package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"push_article/internal/middlewares"
	"push_article/pkg/token"

	"firebase.google.com/go/messaging"
	"github.com/go-chi/chi"
)

type TokenService struct {
	token.Storage
	*messaging.Client
}

func (uts *TokenService) saveToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request struct {
		Token token.Token
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
		return
	}

	userID := middlewares.UserFromContext(r.Context())

	log.Printf("Saving token %s for user %d", request.Token.Token, userID)
	err = uts.SaveToken(r.Context(), userID, request.Token)
	if err != nil {
		http.Error(w, "token save error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (uts *TokenService) userTokens(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "bad or missing user id: "+err.Error(), http.StatusBadRequest)
		return
	}

	tokens, err := uts.UserTokens(r.Context(), userID)
	if err != nil {
		http.Error(w, "token delete error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (uts *TokenService) AddToRouter(r chi.Router) {
	r.With(middlewares.Authorized).Post("/", uts.saveToken)
	r.Get("/", uts.userTokens)
}
