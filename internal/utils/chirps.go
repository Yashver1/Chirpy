package utils

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"unicode/utf8"
)

type Chirp struct {
	Id int `json:"id"`
	Body string `json:"body"`
}

func CreateChirpHandler(database *DB) http.Handler{
	
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		var chirp Chirp
	
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&chirp)
		if err!=nil {
			respondWithErr(w,500,"couldn't read request")
			return
		}
	
		///validation///
	
		if utf8.RuneCountInString(chirp.Body) > 140{
			respondWithErr(w,400,"Chirp is too long")
			return
		}
	
		jsonArray := strings.Split(chirp.Body, " ")
		for i := 0 ; i< len(jsonArray) ; i++{
			lowerString := strings.ToLower(jsonArray[i])
			if lowerString == "kerfuffle" || lowerString == "sharbert" || lowerString == "fornax"{
				jsonArray[i] = "****"
			}
		}
	
		resp := strings.Join(jsonArray, " ")
		///validationEnd///


		createdChirp, err := database.CreateChirp(resp)
		if err!=nil{
			respondWithErr(w,500,"unable to create Chirp")
			return
		}

		respondWithJson(w,201,createdChirp)
	})
}


func GetAllChirpsHandler(database *DB)http.Handler{

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	chirps, err := database.GetChirps()
	if err!=nil{
		respondWithErr(w,500,"unable to get chirps")
		return
	}

	sort.Slice(chirps,func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	respondWithJson(w,200,chirps)

	})
}


