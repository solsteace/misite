package controller

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/solsteace/misite/internal/component/page"
	"github.com/solsteace/misite/internal/persistence"
	"github.com/solsteace/misite/internal/utility/api"
)

const (
	reTagOpPrefix   = `tag:`
	reSerieOpPrefix = `serie:`
	reTitleOpPrefix = `title:`

	reTagOp   = reTagOpPrefix + `[\w,]+`
	reSerieOp = reSerieOpPrefix + `[\w,]+`
	reTitleOp = reTitleOpPrefix + `[\w]+`
	reAnyOp   = `\w`

	reArticle = reTagOp + "|" + reSerieOp + "|" + reAnyOp
	reProject = reTagOp + "|" + reSerieOp + "|" + reAnyOp
	reSerie   = reTitleOp + "|" + reAnyOp
)

func (c Controller) ArticleList(w http.ResponseWriter, r *http.Request) error {
	currentURL, err := url.Parse(r.Header.Get("Hx-Current-URL"))
	if err != nil {
		return fmt.Errorf("controller.ArticleList: %w", err)
	}
	urlQuery := r.URL.Query()

	searchQuery := strings.ToLower(urlQuery.Get("search"))
	lastItem := urlQuery.Get("last")
	param := persistence.ArticlesQueryParam{Last: lastItem}
	if searchQuery != "" {
		if sLimit := urlQuery.Get("limit"); sLimit != "" {
			nLimit, err := strconv.ParseInt(sLimit, 10, strconv.IntSize)
			if err != nil {
				return fmt.Errorf("controller.ArticleList: %w", err)
			} else if nLimit < 0 {
				nLimit = dEFAULT_PAGE_SIZE
			}
			param.Limit = int(nLimit)
		}

		re, err := regexp.Compile(fmt.Sprintf("(%s)+", reArticle))
		if err != nil {
			return fmt.Errorf("controller.ArticleList: %w", err)
		}
		for _, t := range re.FindAll([]byte(searchQuery), -1) {
			token := string(t)
			if p := reTagOpPrefix; strings.HasPrefix(token, p) {
				value := strings.ReplaceAll(strings.TrimPrefix(token, p), "_", " ")
				for _, v := range strings.Split(value, ",") {
					param.Tag = append(param.Tag, v)
				}
			} else if p := reSerieOpPrefix; strings.HasPrefix(token, p) {
				value := strings.ReplaceAll(strings.TrimPrefix(token, p), "_", " ")
				for _, v := range strings.Split(value, ",") {
					param.Serie = append(param.Serie, v)
				}
			} else {
				// Deal with non-operator thing
			}
		}
	}

	articles, err := c.service.Articles(param)
	if err != nil {
		return fmt.Errorf("controller.ArticleList: %w", err)
	}

	var pageComponent templ.Component
	shouldFullRender := (currentURL.Path != api.ExploreArticleUrl || // from outside of the page
		currentURL.Path == api.ExploreArticleUrl && c.isAppRequest(r) && searchQuery == "" && lastItem == "") // calling self via navbar
	if shouldFullRender {
		pageComponent = page.ArticleList(articles)
	} else {
		pageComponent = page.Articles(articles)
	}

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
	currentURL, err := url.Parse(r.Header.Get("Hx-Current-URL"))
	if err != nil {
		return fmt.Errorf("controller.ArticleList: %w", err)
	}
	urlQuery := r.URL.Query()

	searchQuery := strings.ToLower(urlQuery.Get("search"))
	lastItem := urlQuery.Get("last")
	param := persistence.ProjectsQueryParam{Last: lastItem}
	if searchQuery != "" {
		if sLimit := urlQuery.Get("limit"); sLimit != "" {
			nLimit, err := strconv.ParseInt(sLimit, 10, strconv.IntSize)
			if err != nil {
				return fmt.Errorf("controller.ProjectList: %w", err)
			} else if nLimit < 0 {
				nLimit = dEFAULT_PAGE_SIZE
			}
			param.Limit = int(nLimit)
		}

		re, err := regexp.Compile(fmt.Sprintf("(%s)+", reProject))
		if err != nil {
			return fmt.Errorf("controller.ProjectList: %w", err)
		}
		for _, t := range re.FindAll([]byte(searchQuery), -1) {
			token := string(t)
			if p := reTagOpPrefix; strings.HasPrefix(token, p) {
				value := strings.ReplaceAll(strings.TrimPrefix(token, p), "_", " ")
				for _, v := range strings.Split(value, ",") {
					param.Tag = append(param.Tag, v)
				}
			} else if p := reSerieOpPrefix; strings.HasPrefix(token, p) {
				value := strings.ReplaceAll(strings.TrimPrefix(token, p), "_", " ")
				for _, v := range strings.Split(value, ",") {
					param.Serie = append(param.Serie, v)
				}
			} else {
				// Deal with non-operator thing
			}
		}
	}

	projects, err := c.service.Projects(param)
	if err != nil {
		return fmt.Errorf("controller.ProjectList: %w", err)
	}

	var pageComponent templ.Component
	shouldFullRender := (currentURL.Path != api.ExploreProjectUrl || // from outside of the page
		currentURL.Path == api.ExploreProjectUrl && c.isAppRequest(r) && searchQuery == "" && lastItem == "") // calling self via navbar
	if shouldFullRender {
		pageComponent = page.ProjectList(projects)
	} else {
		pageComponent = page.Projects(projects)
	}

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
	currentURL, err := url.Parse(r.Header.Get("Hx-Current-URL"))
	if err != nil {
		return fmt.Errorf("controller.ArticleList: %w", err)
	}
	urlQuery := r.URL.Query()

	searchQuery := strings.ToLower(urlQuery.Get("search"))
	lastItem := urlQuery.Get("last")
	param := persistence.SerieListQueryParam{Last: lastItem}
	if searchQuery != "" {
		if sLimit := urlQuery.Get("limit"); sLimit != "" {
			nLimit, err := strconv.ParseInt(sLimit, 10, strconv.IntSize)
			if err != nil {
				return fmt.Errorf("controller.ArticleList: %w", err)
			} else if nLimit < 0 {
				nLimit = dEFAULT_PAGE_SIZE
			}
			param.Limit = int(nLimit)
		}

		re, err := regexp.Compile(fmt.Sprintf("(%s)+", reSerie))
		if err != nil {
			return fmt.Errorf("controller.ArticleList: %w", err)
		}
		for _, t := range re.FindAll([]byte(searchQuery), -1) {
			token := string(t)
			if p := reTitleOpPrefix; strings.HasPrefix(token, p) {
				value := strings.ReplaceAll(strings.TrimPrefix(token, p), "_", " ")
				param.Title = value
			} else {
				// Deal with non-operator thing
			}
		}
	}

	serieList, err := c.service.SerieList(param)
	if err != nil {
		return fmt.Errorf("controller<Controller.SerieList>: %w", err)
	}

	var pageComponent templ.Component
	shouldFullRender := (currentURL.Path != api.ExploreSeriesUrl || // from outside of the page
		currentURL.Path == api.ExploreSeriesUrl && c.isAppRequest(r) && searchQuery == "" && lastItem == "") // calling self via navbar
	if shouldFullRender {
		pageComponent = page.SerieList(serieList)
	} else {
		pageComponent = page.Series(serieList)
	}

	if !c.isAppRequest(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.SerieList: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.SerieList: %w", err)
	}
	return nil
}

func (c Controller) TagList(w http.ResponseWriter, r *http.Request) error {
	urlQuery := r.URL.Query()
	by := urlQuery.Get("by")
	param := persistence.TagQueryParams{}
	if sLimit := urlQuery.Get("limit"); sLimit != "" {
		nLimit, err := strconv.ParseInt(sLimit, 10, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("controller.TagList: %w", err)
		} else if nLimit < 0 {
			nLimit = dEFAULT_PAGE_SIZE
		}
		param.Limit = int(nLimit)
	}
	if sLimit := urlQuery.Get("page"); sLimit != "" {
		nLimit, err := strconv.ParseInt(sLimit, 10, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("controller.TagList: %w", err)
		} else if nLimit < 0 {
			nLimit = dEFAULT_PAGE_SIZE
		}
		param.Limit = int(nLimit)
	}

	tagStats, err := c.service.Tags(by, param)
	if err != nil {
		return fmt.Errorf("controller.TagList: %w", err)
	}

	pageComponent := page.Tags(by, tagStats)
	if !c.isAppRequest(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.SerieList: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.SerieList: %w", err)
	}
	return nil
}
