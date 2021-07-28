package lib

import "go/ast"

type ParseRequest struct {
	Path       string `json:"path"`
	ProjectKey string `json:"projectKey"`
}

type FunctionCallGraph struct {
	RootPath string         `json:"rootPath"`
	Roots    []FunctionNode `json:"roots"`
}

type Log struct {
	Type   string `json:"type"`
	LogMsg string `json:"log_msg"`
}

type FunctionNode struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Package      string `json:"package"`
	FilePath     string `json:"filePath"`
	Logs         []*Log `json:"logs"`
	ChildNodeIDs []int  `json:"childNodeIDs"`

	Issues []SonarIssue `json:"issues"`

	// for internal use
	funcDecl *ast.FuncDecl
}

type SonarResult struct {
	Total  int
	Issues []SonarIssue
}

type SonarIssue struct {
	Component string
	Line      int
	Severity  string
	Rule      string
	Type      string
	Message   string
	Effort    string
	Debt      string

	// resolve separately
	FilePath string
	Function string
}
