package route

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/solsteace/misite/internal/controller"
	"github.com/solsteace/misite/internal/utility/lib/oops/adapter"
)

type httpHandlerWithError = func(w http.ResponseWriter, r *http.Request) error

type Router struct {
	handler *controller.Controller
}

func NewRouter(handler controller.Controller) Router {
	return Router{handler: &handler}
}

// Kudos: https://boldlygo.tech/posts/2024-01-08-error-handling/
func (r Router) Handle(fx httpHandlerWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := fx(w, req); err != nil {
			ctx := req.Context()
			switch statusCode := adapter.HttpStatusCode(err); statusCode {
			case http.StatusNotFound:
				ctx = context.WithValue(ctx, "err", statusCode)
			default:
				ctx = context.WithValue(ctx, "err", http.StatusInternalServerError)
			}
			ctx = context.WithValue(ctx, "msg",
				fmt.Sprintf("RequestId: %s", middleware.GetReqID(ctx)))

			log.Println(err)
			r.handler.Error(w, req.WithContext(ctx))
		}
	}
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

	router.Get("/project/{id}", r.Handle(r.handler.Project))
	router.Get("/article/{id}", r.Handle(r.handler.Article))
	router.Get("/serie/{id}", r.Handle(r.handler.Serie))
	router.Get("/write", r.Handle(r.handler.MockSpace))
	router.Get("/tags", r.Handle(r.handler.TagList))
	router.Get("/series", r.Handle(r.handler.SerieList))
	router.Get("/articles", r.Handle(r.handler.ArticleList))
	router.Get("/projects", r.Handle(r.handler.ProjectList))
	router.Get("/home", r.Handle(r.handler.Home))
	router.Get("/", r.Handle(r.handler.Home))
	router.NotFound(r.Handle(r.handler.Error))
	parent.Mount("/", router)
}
