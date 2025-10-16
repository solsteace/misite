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

func (c Controller) pageFromName(name string) templ.Component {
	switch name {
	case "home":
		return component.Home()
	case "park":
		return component.Park()
	default:
		return page.NotFound()
	}
}

func (c Controller) ServeBase(w http.ResponseWriter, r *http.Request) error {
	initialBody := "home"
	if initialPage, ok := r.Header["X-Init-Page"]; ok {
		initialBody = initialPage[0]
	}

	ctx := templ.WithChildren(context.Background(), c.pageFromName(initialBody))
	component.
		Base(c.alpinejsUrl, c.htmxUrl).
		Render(ctx, w)
	return nil
}

func (c Controller) ServePage(w http.ResponseWriter, r *http.Request) error {
	name := chi.URLParam(r, "page")

	// Since this webapp uses HTMX, the request with Hx-Request means that
	// the base had already been loaded as it uses HTMX to request the HTML
	// fragment it needs
	//
	// Not fully fool-proof as it could be "spoofed", but this will do
	if _, isHtmx := r.Header["Hx-Request"]; !isHtmx {
		r.Header.Add("X-Init-Page", name)
		return c.ServeBase(w, r)
	}

	c.
		pageFromName(name).
		Render(context.Background(), w)
	return nil
}

func (c Controller) NotFound(w http.ResponseWriter, r *http.Request) error {
	page.NotFound().Render(context.Background(), w)
	w.WriteHeader(http.StatusNotFound)
	return nil
}
