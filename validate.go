package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Error error `json:"error"`
}

type ValidReponse struct {
	Valid bool `json:"valid"`
}

type Chirps struct {
	Body         string `json:"body"`
	Cleaned_Body string `json:"cleaned_body"`
}

func (c *Chirps) getLen() int {
	return len(c.Body)
}

func (c *Chirps) validProfane() {
	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	splitBody := strings.Fields(c.Body)
	for i := 0; i < len(splitBody); i++ {
		for j := 0; j < len(profanity); j++ {
			if strings.EqualFold(splitBody[i], profanity[j]) {
				splitBody[i] = "****"
			}
		}
	}

	c.Cleaned_Body = strings.Join(splitBody, " ")
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

	// clean out profanity
	chirp.validProfane()

	// Checking data

	if chirp.getLen() > 140 {
		jsonError := ErrorResponse{
			Error: err,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(jsonError)
		return
	}

	json.NewEncoder(w).Encode(Chirps{
		Cleaned_Body: chirp.Cleaned_Body,
	})
}
