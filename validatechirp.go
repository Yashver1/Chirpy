package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"
)

type Chirp struct {
	Body string `json:"body"`
}


func respondWithJson( w http.ResponseWriter, statusCode int, payload interface{}) error{
	resp, err := json.Marshal(payload)
	if err!=nil{
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)
	return nil

}

func respondWithErr(w http.ResponseWriter, statusCode int, message string) error {
	if err:=respondWithJson( w, statusCode, map[string]string{"error":message}); err!=nil{
		return err
	}
	return nil
}

func respondWithOk(w http.ResponseWriter) error {
	return respondWithJson(w, 200, map[string]bool{"valid":true})
}


func chirpValidateHandler(w http.ResponseWriter, req *http.Request){
	defer req.Body.Close()
	var chirp Chirp

	type cleanChirp struct{
		Cleaned_body string `json:"cleaned_body"`
	}

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

	cleanedResponse := strings.Join(jsonArray, " ")
	cleanchirp := cleanChirp{
		Cleaned_body: cleanedResponse,
	}
	respondWithJson(w,200, cleanchirp)


	
}