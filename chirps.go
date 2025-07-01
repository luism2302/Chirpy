package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/luism2302/Chirpy/internal/database"
)

type Chirp struct {
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Body       string    `json:"body"`
	User_id    uuid.UUID `json:"user_id"`
}

func cleanBadWords(msg string) string {
	words := strings.Split(msg, " ")
	cleanWords := []string{}
	for _, word := range words {
		lower := strings.ToLower(word)
		if lower == "kerfuffle" || lower == "sharbert" || lower == "fornax" {
			cleanWords = append(cleanWords, "****")
		} else {
			cleanWords = append(cleanWords, word)
		}
	}
	cleanMsg := strings.Join(cleanWords, " ")
	return cleanMsg
}

func respondJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func respondError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	respondJson(w, code, ErrorResponse{Error: msg})
}

func (cfg *apiConfig) chirpsHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	cleanBody := cleanBadWords(params.Body)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldnt decode parameters", err)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleanBody, UserID: params.User_id})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error creating chirp", err)
		return
	}
	respondJson(w, http.StatusCreated, Chirp{Id: chirp.ID, Created_at: chirp.CreatedAt, Updated_at: chirp.UpdatedAt, Body: chirp.Body, User_id: params.User_id})
}
