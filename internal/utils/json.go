package utils

import (
	"encoding/json"
	"net/http"
)


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
