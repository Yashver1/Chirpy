package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"
)

type Chirp struct {
	Id int `json:"id"`
	Body string `json:"body"`
}

func chirpValidateHandler(w http.ResponseWriter, req *http.Request, id int){
	defer req.Body.Close()
	var chirp Chirp

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirp)
	if err!=nil {
		respondWithErr(w,500,"couldn't read request")
		return
	}

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
	chirpCleaned := Chirp{
		Body: resp,
	}
	respondWithJson(w,200, chirpCleaned)
	
}