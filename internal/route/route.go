package route

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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
			log.Println(err)
			statusCode := adapter.HttpStatusCode(err)
			switch statusCode {
			case http.StatusNotFound:
				r.handler.NotFound(w, req)
			default:
				msg := adapter.HttpErrorMsg(err)
				payload := map[string]any{"msg": msg}
				resPayload, err := json.Marshal(payload)
				if err != nil {
					resPayload = []byte("Something went wrong in our system")
				}

				w.WriteHeader(statusCode)
				w.Write(resPayload)
			}
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
	router.NotFound(r.Handle(r.handler.NotFound))
	parent.Mount("/", router)
}
