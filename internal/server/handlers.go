package server

import (
	"html/template"
	"net/http"

	"github.com/d-jackalope/L0/pkg/models"
)

type templateHome struct {
	Lenght int
	Orders []models.Order
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		s.notFound(w)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		s.serverError(w, err)
		return
	}
	orders := s.App.Cache.RandomOrders()
	lenght := s.App.Cache.CacheSize()
	data := &templateHome{Orders: orders, Lenght: lenght}
	err = ts.Execute(w, data)
	if err != nil {
		s.serverError(w, err)
		return
	}
}

func (s *Server) order(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Query().Get("uid")
	order, exist := s.App.Cache.Get(orderUID)
	if !exist {
		s.notFound(w)
		return
	}

	files := []string{
		"./ui/html/order.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		s.serverError(w, err)
		return
	}

	err = ts.Execute(w, order)
	if err != nil {
		s.serverError(w, err)
		return
	}
}

func (s *Server) notFound(w http.ResponseWriter) {
	files := []string{
		"./ui/html/notfound.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		s.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "notfound", nil)
	if err != nil {
		s.serverError(w, err)
		return
	}

}
