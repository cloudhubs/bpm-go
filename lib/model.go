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

type Log struct {
	ID          int        `json:"id"`
	LineNumber  int        `json:"line"`
	LogMsg      string     `json:"log_msg"`
}

type FunctionNode struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Package      string   `json:"package"`
	FilePath     string   `json:"filePath"`
	Logs         []Log 	  `json:"logs"`
	ChildNodeIDs []int    `json:"childNodeIDs"`

	// for internal use
	funcDecl *ast.FuncDecl
}
