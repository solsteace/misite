package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/solsteace/misite/internal/component/page"
	"github.com/solsteace/misite/internal/component/widget"
	"github.com/solsteace/misite/internal/persistence"
)

func (c Controller) ArticleList(w http.ResponseWriter, r *http.Request) error {
	urlQuery := r.URL.Query()
	param := persistence.ArticlesQueryParam{}
	if c.isAppRequest(r) {
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
	if !c.isAppRequest(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.ArticleList: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.ArticleList: %w", err)
	}
	return nil
}

func (c Controller) ProjectList(w http.ResponseWriter, r *http.Request) error {
	urlQuery := r.URL.Query()
	param := persistence.ProjectsQueryParam{}
	if c.isAppRequest(r) {
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
	if !c.isAppRequest(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.ProjectList: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.ProjectList: %w", err)
	}
	return nil
}

func (c Controller) SerieList(w http.ResponseWriter, r *http.Request) error {
	urlQuery := r.URL.Query()
	param := persistence.SerieListQueryParam{}
	if c.isAppRequest(r) {
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
	}

	serieList, err := c.service.SerieList(param)
	if err != nil {
		return fmt.Errorf("controller<Controller.SerieList>: %w", err)
	}

	pageComponent := page.SerieList(serieList)
	if !c.isAppRequest(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.SerieList: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.SerieList: %w", err)
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
