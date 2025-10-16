package controller

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/solsteace/misite/internal/component"
	"github.com/solsteace/misite/internal/component/page"
	"github.com/solsteace/misite/internal/service"
)

type Controller struct {
	service service.Service

	alpinejsUrl string // url to alpinejs script (unrelated to controller, but we're gonna stick with these infra anyway for now)
	htmxUrl     string // url to htmx script (unrelated to controller, but we're gonna stick with these infra anyway for now)
}

func NewController(
	service service.Service,
	alpinejsUrl string,
	htmxUrl string,
) Controller {
	return Controller{
		service:     service,
		alpinejsUrl: alpinejsUrl,
		htmxUrl:     htmxUrl}
}

func (c Controller) NotFound(w http.ResponseWriter, r *http.Request) error {
	ctx := templ.WithChildren(context.Background(), page.NotFound())
	component.Base(c.alpinejsUrl, c.htmxUrl).Render(ctx, w)

	w.WriteHeader(http.StatusNotFound)
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

	ctx := templ.WithChildren(context.Background(), child)
	component.Base(c.alpinejsUrl, c.htmxUrl).Render(ctx, w)
	return nil
}
