package lib

import (
	"go/parser"
	"go/token"
)

func GetImportContext(request ParseRequest) (ImportContext, error) {
	importContext := ImportContext{
		RootPath: request.Path,
	}

	sources, err := WalkMatch(request.Path, "*.go")
	if err != nil {
		return ImportContext{}, err
	}

	for _, source := range sources {
		file, err := parser.ParseFile(token.NewFileSet(), source, nil, parser.ImportsOnly)
		if err != nil {
			return ImportContext{}, err
		}

		importForFile := ImportForFile{
			Path: source,
		}

		for _, item := range file.Imports {
			importForFile.Packages = append(importForFile.Packages, RemoveQuotes(item.Path.Value))
		}

		importContext.Imports = append(importContext.Imports, importForFile)
	}

	return importContext, err
}
