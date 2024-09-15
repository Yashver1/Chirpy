package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"unicode/utf8"
	
	"github.com/Yashver1/chirpy/internal/utils"
	
)



func (a *apiConfig) CreateChirpHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		var chirp utils.Chirp
		database := a.Database

		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&chirp)
		if err != nil {
			
			utils.RespondWithErr(w, 500, "couldn't read request")
			return
		}

		///validation///

		if utf8.RuneCountInString(chirp.Body) > 140 {
			utils.RespondWithErr(w, 400, "Chirp is too long")
			return
		}

		jsonArray := strings.Split(chirp.Body, " ")
		for i := 0; i < len(jsonArray); i++ {
			lowerString := strings.ToLower(jsonArray[i])
			if lowerString == "kerfuffle" || lowerString == "sharbert" || lowerString == "fornax" {
				jsonArray[i] = "****"
			}
		}

		resp := strings.Join(jsonArray, " ")
		///validationEnd///

		createdChirp, err := database.CreateChirp(resp)
		if err != nil {
			utils.RespondWithErr(w, 500, "unable to create Chirp")
			return
		}

		utils.RespondWithJSON(w, 201, createdChirp)
	})
}

func (a *apiConfig) GetAllChirpsHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		database := a.Database
		chirps, err := database.GetChirps()
		if err != nil {
			utils.RespondWithErr(w, 500, "unable to get chirps")
			return
		}

		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id < chirps[j].Id
		})

		utils.RespondWithJSON(w, 200, chirps)

	})
}

func (a *apiConfig) GetChirpHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		database := a.Database
		idString := r.PathValue("chirpID")
		chirp, err := database.GetChirp(idString)
		if err != nil {
			utils.RespondWithErr(w, 404, "Chirp does not exist")
			return
		}

		utils.RespondWithJSON(w, 200, chirp)
	})
}
