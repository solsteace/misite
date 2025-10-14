package controller

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/solsteace/misite/internal/component"
	"github.com/solsteace/misite/internal/component/page"
	"github.com/solsteace/misite/internal/service"
)

type Controller struct {
	service service.Service
}

func (c Controller) NotFound(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNotFound)
	component.RenderBaseTo(w, page.NotFound())
	return nil
}

func (c Controller) RedirectToHome(w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
	return nil
}

func (c Controller) ServePage(w http.ResponseWriter, r *http.Request) error {
	var child templ.Component
	switch chi.URLParam(r, "page") {
	case "home":
		child = component.Home()
	case "park":
		child = component.Park()
	default:
		return c.NotFound(w, r)
	}

	component.RenderBaseTo(w, child)
	return nil
}

func NewController(service service.Service) Controller {
	return Controller{service: service}
}
