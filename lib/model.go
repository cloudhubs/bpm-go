package lib

import "go/ast"

type ParseRequest struct {
	Path string `json:"path"`
}

type ImportContext struct {
	RootPath string          `json:"rootPath"`
	Imports  []ImportForFile `json:"imports"`
}

type ImportForFile struct {
	Path     string   `json:"path"`
	Packages []string `json:"packages"`
}

type FunctionCallGraph struct {
	RootPath string         `json:"rootPath"`
	Roots    []FunctionNode `json:"roots"`
}

type FunctionNode struct {
	Name       string         `json:"name"`
	Package    string         `json:"package"`
	FilePath   string         `json:"filePath"`
	Logs       []string       `json:"logs"`
	ChildNodes []FunctionNode `json:"childNodes"`

	// for internal use
	funcDecl *ast.FuncDecl
}
