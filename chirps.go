package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type ChirpRequest struct {
	Body string `json:"body"`
}

type ErrorJson struct {
	Error string `json:"error"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {

	dbChirps, err := cfg.DB.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirpRequest := ChirpRequest{}
	err := decoder.Decode(&chirpRequest)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// check for validity
	if len(chirpRequest.Body) > 140 {
		// throw error
		err1 := errors.New("Chirp is too long")
		respondWithError(w, http.StatusBadRequest, err1.Error())
		return
	}

	// normal flow
	ChirpJsonResponse := Chirp{}

	// Replacement logic
	{
		// Input string
		input := chirpRequest.Body

		// Words to be replaced
		wordsToReplace := []string{"kerfuffle", "sharbert", "fornax"}

		// Case-insensitive regular expression pattern
		// The \b word boundary ensures that only full words are matched
		pattern := "(?i)\\b(" + strings.Join(wordsToReplace, "|") + ")\\b"

		// Create the regular expression
		re := regexp.MustCompile(pattern)

		// Replace the words with "****"
		ChirpJsonResponse.Body = re.ReplaceAllString(input, "****")
	}

	chirp, err := cfg.DB.CreateChirp(ChirpJsonResponse.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:   chirp.ID,
		Body: chirp.Body,
	})
}
