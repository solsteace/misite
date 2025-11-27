package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/solsteace/misite/internal/component"
	"github.com/solsteace/misite/internal/component/page"
	"github.com/solsteace/misite/internal/service"
)

const (
	tagQueryParam   = "tId"
	serieQueryParam = "sId"

	dEFAULT_PAGE_SIZE = 10
)

type Controller struct {
	service service.Service

	indexUrl    string // url to homepage
	alpinejsUrl string // url to alpinejs script (unrelated to controller, but we're gonna stick with these infra anyway for now)
	htmxUrl     string // url to htmx script (unrelated to controller, but we're gonna stick with these infra anyway for now)
}

func NewController(
	service service.Service,
	indexUrl string,
	alpinejsUrl string,
	htmxUrl string,
) Controller {
	return Controller{
		service:     service,
		indexUrl:    indexUrl,
		alpinejsUrl: alpinejsUrl,
		htmxUrl:     htmxUrl}
}

// This is not totally fool-proof as it could be "spoofed". Better way? maybe next time
func (c Controller) isAppRequest(r *http.Request) bool {
	_, ok := r.Header["Hx-Request"]
	return ok
}

// Serves a page with its base
func (c Controller) serveWithBase(
	body templ.Component,
	w http.ResponseWriter,
	r *http.Request,
) error {
	ctx := templ.WithChildren(context.Background(), body)
	if err := component.Base(c.alpinejsUrl, c.htmxUrl).Render(ctx, w); err != nil {
		return fmt.Errorf("controller.serveWithBase: %w", err)
	}
	return nil
}

func (c Controller) Home(w http.ResponseWriter, r *http.Request) error {
	pageComponent := page.Home(c.indexUrl)
	if !c.isAppRequest(r) {
		return c.serveWithBase(pageComponent, w, r)
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.Home: %w", err)
	}
	return nil
}

func (c Controller) Error(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	extraMesssage, ok := ctx.Value("msg").(string)
	if !ok {
		extraMesssage = "additinal info missing, sorry"
	}
	code, ok := ctx.Value("err").(int)
	if !ok {
		code = http.StatusInternalServerError
	}
	pageComponent := page.Error(code, extraMesssage)
	if !c.isAppRequest(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.Error: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.Error: %w", err)
	}
	return nil
}
