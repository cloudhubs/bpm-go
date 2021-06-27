package api

import (
	"bpm-go/lib"
	"encoding/json"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	respondMessage(w, http.StatusOK, "hello")
}

func getImports(w http.ResponseWriter, r *http.Request) {
	request := lib.ParseRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	importContext, err := lib.GetImportContext(request)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	} else {
		respondJSON(w, http.StatusOK, importContext)
	}
}
