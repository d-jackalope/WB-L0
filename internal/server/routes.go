package server

import "net/http"

func (s *Server) routers() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.home)
	mux.HandleFunc("/order", s.order)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
