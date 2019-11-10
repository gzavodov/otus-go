package web

import (
	"fmt"
	"net/http"
	"path"

	"go.uber.org/zap"
)

type Server struct {
	Logger     *zap.Logger
	HTTPServer *http.ServeMux
}

func (h *Server) Start(address string) error {
	h.HTTPServer = http.NewServeMux()
	h.HTTPServer.Handle("/", h)
	return http.ListenAndServe(address, h.HTTPServer)
}

func (h Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Request handling", zap.String("URL Path", r.URL.Path))
	isMatch, _ := path.Match("/hello/*", r.URL.Path)
	if !isMatch {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Welcome to the home page!")
}