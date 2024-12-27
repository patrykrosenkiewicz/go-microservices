package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Port   string
	server *http.Server
	router *mux.Router
}

// handleIndex handles requests to the root ("/") and prints a personalized message.
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hi there, I talk from notification %s!", r.URL.Path[1:])
}

// NewServer creates and returns a new Server instance.
func NewServer(port string) *Server {
	if port == "" {
		port = "8080"
	}

	s := &Server{
		Port:   port,
		server: &http.Server{},
		router: mux.NewRouter(),
	}

	// Define routes
	router := s.router.PathPrefix("/notification").Subrouter()
	router.HandleFunc("/{name}", s.handleIndex).Methods("GET")

	return s
}

// Open starts the server on the specified port.
func (s *Server) Open() error {
	s.server.Addr = ":" + s.Port
	s.server.Handler = s.router
	return s.server.ListenAndServe()
}
