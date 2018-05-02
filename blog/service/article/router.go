package article

import (
	"github.com/gorilla/mux"
)

func (s *Service) InitRouter() {
	r := mux.NewRouter().PathPrefix("/api/blog.Article").Subrouter()

	r.HandleFunc("/Hello", s.GenHTTPHandler(s.Hello)).Methods("POST")

	s.RegisterRouter(r)
}
