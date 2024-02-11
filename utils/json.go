package utils

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, responseObject interface{}) error {
	return RespondWithJSON(w, http.StatusOK, responseObject)
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, responseObject interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(responseObject)
}
