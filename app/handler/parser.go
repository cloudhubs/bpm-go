package handler

import (
	"encoding/json"
	"net/http"

	"rad-go/app/model"
)

func ListImports(w http.ResponseWriter, r *http.Request) {
	request := model.ParseRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	respondJSON(w, http.StatusOK, request)
}
