package utils

import (
	"encoding/json"
	"net/http"
)


func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) error {
	resp, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)
	return nil
}

func RespondWithErr(w http.ResponseWriter, statusCode int, message string) error {
	if err := RespondWithJSON(w, statusCode, map[string]string{"error": message}); err != nil {
		return err
	}
	return nil
}

func RespondWithOk(w http.ResponseWriter) error {
	return RespondWithJSON(w, 200, map[string]bool{"valid": true})
}
