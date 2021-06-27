package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// respondJSON makes the response with payload as json format.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, errWrite := w.Write([]byte(err.Error())); errWrite != nil {
			log.Println(errWrite)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, errWrite := w.Write(response); errWrite != nil {
		log.Println(errWrite)
	}
}

// respondMessage attaches a message with the response payload as json format.
func respondMessage(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"message": message})
}

// respondError makes the error response with payload as json format.
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
