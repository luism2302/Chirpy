package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldnt decode parameters", err)
	}
	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error creating user", err)
		return
	}
	respondJson(w, http.StatusCreated, User{Id: user.ID, Created_at: user.CreatedAt, Updated_at: user.CreatedAt, Email: user.Email})
}
