package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/solsteace/misite/internal/component"
	"github.com/solsteace/misite/internal/component/page"
	"github.com/solsteace/misite/internal/component/widget"
	"github.com/solsteace/misite/internal/persistence"
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
	if err := component.Base(c.alpinejsUrl, c.htmxUrl).Render(ctx, w); err != nil {
		return fmt.Errorf("controller.serveWithBase: %w", err)
	}
	return nil
}

func (c Controller) Home(w http.ResponseWriter, r *http.Request) error {
	pageComponent := page.Home(c.indexUrl)
	if c.requestNeedsBase(r) {
		return c.serveWithBase(pageComponent, w, r)
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.Home: %w", err)
	}
	return nil
}

func (c Controller) ArticleList(w http.ResponseWriter, r *http.Request) error {
	urlQuery := r.URL.Query()
	param := persistence.ArticlesQueryParam{}
	if sPage := urlQuery.Get("page"); sPage != "" {
		nPage, err := strconv.ParseInt(sPage, 10, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("controller.ArticleList: %w", err)
		} else if nPage < 0 {
			nPage = 0
		}
		param.Page = int(nPage)
	}
	if sLimit := urlQuery.Get("limit"); sLimit != "" {
		nLimit, err := strconv.ParseInt(sLimit, 10, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("controller.ArticleList: %w", err)
		} else if nLimit < 0 {
			nLimit = dEFAULT_PAGE_SIZE
		}
		param.Limit = int(nLimit)
	}
	for _, id := range urlQuery[tagQueryParam] {
		nId, err := strconv.ParseInt(id, 10, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("controller.ArticleList: %w", err)
		} else if nId > 0 {
			param.TagId = append(param.TagId, int(nId))
		}
	}
	for _, id := range urlQuery[serieQueryParam] {
		nId, err := strconv.ParseInt(id, 10, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("controller.ArticleList: %w", err)
		} else if nId > 0 {
			param.SerieId = append(param.SerieId, int(nId))
		}
	}

	articles, err := c.service.Articles(param)
	if err != nil {
		return fmt.Errorf("controller.ArticleList: %w", err)
	}

	var sTopTags string
	topTags := c.service.MostTagsOnArticles(articles, 16)
	for idx, i := range topTags {
		sTopTags += fmt.Sprintf("%s=%s", tagQueryParam, strconv.Itoa(i))
		if idx < len(topTags)-1 {
			sTopTags += "&"
		}
	}

	pageComponent := page.ArticleList(articles, sTopTags)
	if c.requestNeedsBase(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.ArticleList: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.ArticleList: %w", err)
	}
	return nil
}

func (c Controller) Article(w http.ResponseWriter, r *http.Request) error {
	articleId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, strconv.IntSize)
	if err != nil {
		return fmt.Errorf("controller.Article: %w", err)
	}

	article, err := c.service.Article(int(articleId))
	if err != nil {
		return fmt.Errorf("controller.Article: %w", err)
	}

	pageComponent := page.Article(article)
	if c.requestNeedsBase(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.Article: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.Article: %w", err)
	}
	return nil
}

func (c Controller) ProjectList(w http.ResponseWriter, r *http.Request) error {
	urlQuery := r.URL.Query()
	param := persistence.ProjectsQueryParam{}
	if sPage := urlQuery.Get("page"); sPage != "" {
		nPage, err := strconv.ParseInt(sPage, 10, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("controller.ProjectList: %w", err)
		} else if nPage < 0 {
			nPage = 0
		}
		param.Page = int(nPage)
	}
	if sLimit := urlQuery.Get("limit"); sLimit != "" {
		nLimit, err := strconv.ParseInt(sLimit, 10, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("controller.ProjectList: %w", err)
		} else if nLimit < 0 {
			nLimit = dEFAULT_PAGE_SIZE
		}
		param.Limit = int(nLimit)
	}
	for _, id := range urlQuery[tagQueryParam] {
		nId, err := strconv.ParseInt(id, 10, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("controller.ProjectList: %w", err)
		} else if nId > 0 {
			param.TagId = append(param.TagId, int(nId))
		}
	}

	projects, err := c.service.Projects(param)
	if err != nil {
		return fmt.Errorf("controller.ProjectList: %w", err)
	}

	var sTopTags string
	topTags := c.service.MostTagsOnProjects(projects, 16)
	for idx, i := range topTags {
		sTopTags += fmt.Sprintf("%s=%s", tagQueryParam, strconv.Itoa(i))
		if idx < len(topTags)-1 {
			sTopTags += "&"
		}
	}

	pageComponent := page.ProjectList(projects, sTopTags)
	if c.requestNeedsBase(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.ProjectList: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.ProjectList: %w", err)
	}
	return nil
}

func (c Controller) Project(w http.ResponseWriter, r *http.Request) error {
	projectId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, strconv.IntSize)
	if err != nil {
		return fmt.Errorf("controller.Project: %w", err)
	}

	project, err := c.service.Project(int(projectId))
	if err != nil {
		return fmt.Errorf("controller.Project: %w", err)
	}

	pageComponent := page.Project(project)
	if c.requestNeedsBase(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.Project: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.Project: %w", err)
	}
	return nil
}

func (c Controller) SerieList(w http.ResponseWriter, r *http.Request) error {
	urlQuery := r.URL.Query()
	param := persistence.SerieListQueryParam{}
	if sPage := urlQuery.Get("page"); sPage != "" {
		nPage, err := strconv.ParseInt(sPage, 10, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("controller.SerieList: %w", err)
		} else if nPage < 0 {
			nPage = 0
		}
		param.Page = int(nPage)
	}
	if sLimit := urlQuery.Get("limit"); sLimit != "" {
		nLimit, err := strconv.ParseInt(sLimit, 10, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("controller.SerieList: %w", err)
		} else if nLimit < 0 {
			nLimit = dEFAULT_PAGE_SIZE
		}
		param.Limit = int(nLimit)
	}

	serieList, err := c.service.SerieList(param)
	if err != nil {
		return fmt.Errorf("controller<Controller.SerieList>: %w", err)
	}

	pageComponent := page.SerieList(serieList)
	if c.requestNeedsBase(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.SerieList: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.SerieList: %w", err)
	}
	return nil
}

func (c Controller) Serie(w http.ResponseWriter, r *http.Request) error {
	serieId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, strconv.IntSize)
	if err != nil {
		return fmt.Errorf("controller.Serie: %w", err)
	}

	serie, err := c.service.Serie(int(serieId))
	if err != nil {
		return fmt.Errorf("controller<Controller.Serie>; %w", err)
	}

	pageComponent := page.Serie(serie)
	if c.requestNeedsBase(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.Serie: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.Serie: %w", err)
	}
	return nil
}

func (c Controller) NotFound(w http.ResponseWriter, r *http.Request) error {
	pageComponent := page.NotFound()
	w.WriteHeader(http.StatusNotFound)
	if c.requestNeedsBase(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.NotFound: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.NotFound: %w", err)
	}
	return nil
}

func (c Controller) ExplorationTags(w http.ResponseWriter, r *http.Request) error {
	urlQuery := r.URL.Query()

	var tagIds []int
	for _, id := range urlQuery[tagQueryParam] {
		nId, err := strconv.ParseInt(id, 10, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("controller.Tags: %w", err)
		} else if nId > 0 {
			tagIds = append(tagIds, int(nId))
		}
	}

	var pageComponent templ.Component
	switch urlQuery.Get("by") {
	case "article":
		tags, count, err := c.service.CountArticleMatchingTags(tagIds)
		if err != nil {
			return fmt.Errorf("controller<Controller.Tags>: %w", err)
		}
		pageComponent = widget.ExplorationTags(tags, count, "articles")
	case "project":
		tags, count, err := c.service.CountProjectMatchingTags(tagIds)
		if err != nil {
			return fmt.Errorf("controller<Controller.Tags>: %w", err)
		}
		pageComponent = widget.ExplorationTags(tags, count, "projects")
	default:
		pageComponent = page.NotFound()
	}
	if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.Tags: %w", err)
	}
	return nil
}
