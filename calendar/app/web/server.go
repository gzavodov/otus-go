package web

import (
	"fmt"
	"net/http"
	"path"

	"go.uber.org/zap"
)

//NewServer Creates new web server
func NewServer(c *Config, logger *zap.Logger) *Server {
	server := &Server{Address: c.HTTPAddress, Logger: logger}
	return server
}
//Server Simple Web Server
type Server struct {
	Address    string
	Logger     *zap.Logger
	HTTPServer *http.ServeMux
}

//Start Start handling of web requests 
func (h *Server) Start() error {
	h.HTTPServer = http.NewServeMux()
	h.HTTPServer.Handle("/", h)
	return http.ListenAndServe(h.Address, h.HTTPServer)
}

//ServeHTTP Web request handler
func (h Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Request handling", zap.String("URL Path", r.URL.Path))
	isMatch, _ := path.Match("/hello/*", r.URL.Path)
	if !isMatch {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Welcome to the home page!")
}
