package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/jackellisw/httpNoFramework.git/internal/database"
)

type Emails struct {
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) CreateUsersHandler(w http.ResponseWriter, r *http.Request, db *database.Queries) {
	// Takes in an email
	var email Emails
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&email); err != nil {
		http.Error(w, "Could not decode email", http.StatusBadRequest)
		return
	}

	// Returns the a new user's id, email, and timestamp
	user, err := cfg.db.CreateUser(r.Context(), email.Email)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     email.Email,
	})
}
