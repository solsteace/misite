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

	router.Get("/project/{id}", reqres.HttpHandlerWithError(r.handler.Project))
	router.Get("/article/{id}", reqres.HttpHandlerWithError(r.handler.Article))
	router.Get("/serie/{id}", reqres.HttpHandlerWithError(r.handler.Serie))
	router.Get("/write", reqres.HttpHandlerWithError(r.handler.MockSpace))
	router.Get("/tags", reqres.HttpHandlerWithError(r.handler.TagList))
	router.Get("/series", reqres.HttpHandlerWithError(r.handler.SerieList))
	router.Get("/articles", reqres.HttpHandlerWithError(r.handler.ArticleList))
	router.Get("/projects", reqres.HttpHandlerWithError(r.handler.ProjectList))
	router.Get("/home", reqres.HttpHandlerWithError(r.handler.Home))
	router.Get("/", reqres.HttpHandlerWithError(r.handler.Home))
	router.NotFound(reqres.HttpHandlerWithError(r.handler.NotFound))
	parent.Mount("/", router)
}
