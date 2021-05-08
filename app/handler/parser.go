package handler

import (
	"encoding/json"
	"go/parser"
	"go/token"
	"net/http"
	"os"
	"path/filepath"

	"rad-go/app/model"
)

func GetImportContext(w http.ResponseWriter, r *http.Request) {
	request := model.ParseRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	importContext := model.ImportContext{
		RootPath: request.Path,
	}

	sources, err := WalkMatch(request.Path, "*.go")
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	for _, source := range sources {
		file, err := parser.ParseFile(token.NewFileSet(), source, nil, parser.ImportsOnly)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		importForFile := model.ImportForFile{
			Path: source,
		}

		for _, item := range file.Imports {
			importForFile.Packages = append(importForFile.Packages, RemoveQuotes(item.Path.Value))
		}

		importContext.Imports = append(importContext.Imports, importForFile)
	}

	respondJSON(w, http.StatusOK, importContext)
}

func RemoveQuotes(str string) string {
	return str[1 : len(str)-1]
}

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
