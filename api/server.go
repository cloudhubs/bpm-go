package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Server has router instance and other configuration to start the API server.
type Server struct {
	Port   int
	router *mux.Router
}

// InitializeAndRun initializes and run the API server on it's router.
func (s *Server) InitializeAndRun() {
	// initialize and set routers
	s.router = mux.NewRouter()
	s.route("/hello", http.MethodGet, hello)
	s.route("/imports", http.MethodPost, getImports)
	s.route("/functions", http.MethodPost, getFunctionNodes)
	s.route("/sonar", http.MethodPost, runSonarAnalysis)

	// start the server
	addr := fmt.Sprintf(":%v", s.Port)
	log.Println("Server started at", addr)
	log.Fatalln(http.ListenAndServe(addr, s.router))
}

// Route wraps the router for a HTTP method.
func (s *Server) route(path string, method string, handler RequestHandlerFunction) {
	s.router.HandleFunc(path, handler).Methods(method)
}

// RequestHandlerFunction defines a common type for request handler.
type RequestHandlerFunction func(w http.ResponseWriter, r *http.Request)
