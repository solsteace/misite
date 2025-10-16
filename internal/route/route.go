package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/solsteace/misite/internal/controller"
	"github.com/solsteace/misite/internal/utility/lib/reqres"
)

type Router struct {
	handler *controller.Controller
}

func NewRouter(handler controller.Controller) Router {
	return Router{handler: &handler}
}

func (r Router) UseOn(parent *chi.Mux) {
	router := chi.NewRouter()

	// Basic file server
	// Might be customized to restrict which files are retrievable
	router.Get(
		"/static/*",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("./static"))).ServeHTTP)

	router.Get("/{page}", reqres.HttpHandlerWithError(r.handler.ServePage))
	router.Get("/", reqres.HttpHandlerWithError(r.handler.ServeBase))
	router.NotFound(reqres.HttpHandlerWithError(r.handler.NotFound))
	parent.Mount("/", router)
}
