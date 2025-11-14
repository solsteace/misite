package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/solsteace/misite/internal/component/page"
	"github.com/solsteace/misite/internal/persistence"
	"github.com/solsteace/misite/internal/utility/lib/oops"
)

// MEGA TODO: REFACTOR
func (c Controller) Search(w http.ResponseWriter, r *http.Request) error {
	if !c.isAppRequest(r) {
		return oops.Forbidden{}
	}

	prompt := r.FormValue("prompt")
	switch urlQuery := r.URL.Query(); urlQuery.Get("cat") {
	case "article":
		if prompt == "" {
			param := persistence.ArticlesQueryParam{}
			if sPage := urlQuery.Get("page"); sPage != "" {
				nPage, err := strconv.ParseInt(sPage, 10, strconv.IntSize)
				if err != nil {
					return fmt.Errorf("controller.Search: %w", err)
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

			articles, err := c.service.Articles(param)
			if err != nil {
				return fmt.Errorf("controller.Search: %w", err)
			}

			pageComponent := page.Articles(articles)
			if !c.isAppRequest(r) {
				if err := c.serveWithBase(pageComponent, w, r); err != nil {
					return fmt.Errorf("controller.Search: %w", err)
				}
			} else if err := pageComponent.Render(context.Background(), w); err != nil {
				return fmt.Errorf("controller.Search: %w", err)
			}
			return nil
		}
	case "project":
		if prompt == "" {
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

			projects, err := c.service.Projects(param)
			if err != nil {
				return fmt.Errorf("controller.ProjectList: %w", err)
			}

			pageComponent := page.Projects(projects)
			if !c.isAppRequest(r) {
				if err := c.serveWithBase(pageComponent, w, r); err != nil {
					return fmt.Errorf("controller.ProjectList: %w", err)
				}
			} else if err := pageComponent.Render(context.Background(), w); err != nil {
				return fmt.Errorf("controller.ProjectList: %w", err)
			}
			return nil
		}
	case "serie":
		if prompt == "" {
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

			pageComponent := page.Series(serieList)
			if !c.isAppRequest(r) {
				if err := c.serveWithBase(pageComponent, w, r); err != nil {
					return fmt.Errorf("controller.SerieList: %w", err)
				}
			} else if err := pageComponent.Render(context.Background(), w); err != nil {
				return fmt.Errorf("controller.SerieList: %w", err)
			}
			return nil
		}
	}
	return oops.BadRequest{}
}
