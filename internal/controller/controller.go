package controller

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
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

// Checks whether the request only needs the page component or the whole full page
//
// Assuming all component only requests came from the app which uses HTMX,
// requests for fragments without `Hx-Request` would be assumed as a request
// for a full page (base + wanted page). This is not totally fool-proof as it
// could be "spoofed", but this will do for now
func (c Controller) requestNeedsBase(r *http.Request) bool {
	_, ok := r.Header["Hx-Request"]
	return !ok
}

// Serves a page with its base
func (c Controller) serveWithBase(
	body templ.Component,
	w http.ResponseWriter,
	r *http.Request,
) error {
	ctx := templ.WithChildren(context.Background(), body)
	component.
		Base(c.alpinejsUrl, c.htmxUrl).
		Render(ctx, w)
	return nil
}

func (c Controller) Home(w http.ResponseWriter, r *http.Request) error {
	pageComponent := page.Home()
	if c.requestNeedsBase(r) {
		return c.serveWithBase(pageComponent, w, r)
	}
	pageComponent.Render(context.Background(), w)
	return nil
}

func (c Controller) Articles(w http.ResponseWriter, r *http.Request) error {
	pageComponent := page.ArticleList()
	if c.requestNeedsBase(r) {
		return c.serveWithBase(pageComponent, w, r)
	}
	pageComponent.Render(context.Background(), w)
	return nil
}

func (c Controller) Article(w http.ResponseWriter, r *http.Request) error {
	pageComponent := page.Article()
	if c.requestNeedsBase(r) {
		return c.serveWithBase(pageComponent, w, r)
	}
	pageComponent.Render(context.Background(), w)
	return nil
}

func (c Controller) Projects(w http.ResponseWriter, r *http.Request) error {
	pageComponent := page.ProjectList()
	if c.requestNeedsBase(r) {
		return c.serveWithBase(pageComponent, w, r)
	}
	pageComponent.Render(context.Background(), w)
	return nil
}

func (c Controller) Project(w http.ResponseWriter, r *http.Request) error {
	pageComponent := page.Project()
	if c.requestNeedsBase(r) {
		return c.serveWithBase(pageComponent, w, r)
	}
	pageComponent.Render(context.Background(), w)
	return nil
}

func (c Controller) NotFound(w http.ResponseWriter, r *http.Request) error {
	pageComponent := page.NotFound()
	if c.requestNeedsBase(r) {
		return c.serveWithBase(pageComponent, w, r)
	}
	pageComponent.Render(context.Background(), w)
	w.WriteHeader(http.StatusNotFound)
	return nil
}
