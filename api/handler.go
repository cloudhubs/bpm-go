package api

import (
	"bpm-go/lib"
	"encoding/json"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	respondMessage(w, http.StatusOK, "hello")
}

func getFunctionNodes(w http.ResponseWriter, r *http.Request) {
	request := lib.ParseRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	fnNodes, err := lib.GetFunctionNodes(request)

	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	} else {
		respondJSON(w, http.StatusOK, fnNodes)
	}
}

func runSonarAnalysis(w http.ResponseWriter, r *http.Request) {
	request := lib.ParseRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	resp, err := lib.RunSonarAnalysis(request)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	} else {
		respondJSON(w, http.StatusOK, resp)
	}
}

func runProjectAnalysis(w http.ResponseWriter, r *http.Request) {
	request := lib.ParseRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	resp, err := lib.RunProjectAnalysis(request)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	} else {
		respondJSON(w, http.StatusOK, resp)
	}
}
