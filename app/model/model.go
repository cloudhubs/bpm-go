package model

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
