package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Error error `json:"error"`
}

type ValidReponse struct {
	Valid bool `json:"valid"`
}

type Chirps struct {
	Body string `json:"body"`
}

func (c *Chirps) getLen() int {
	return len(c.Body)
}

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	// Decode data
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	chirp := Chirps{}
	err := decoder.Decode(&chirp)
	if err != nil {
		jsonError := ErrorResponse{
			Error: fmt.Errorf("Couldn't decode data: %s", err),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(jsonError)
	}

	// Checking data

	if chirp.getLen() > 140 {
		jsonError := ErrorResponse{
			Error: err,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(jsonError)
		return
	}

	// encode reponse with isValid or Error
	validRes := ValidReponse{
		Valid: true,
	}

	json.NewEncoder(w).Encode(validRes)
}
