package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/solsteace/misite/internal/component/page"
	"github.com/solsteace/misite/internal/persistence"
)

func (c Controller) Article(w http.ResponseWriter, r *http.Request) error {
	articleId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, strconv.IntSize)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			articleId = -1
		} else {
			return fmt.Errorf("controller.Article: %w", err)
		}
	}

	article, err := c.service.Article(int(articleId))
	if err != nil {
		return fmt.Errorf("controller.Article: %w", err)
	}

	pageComponent := page.Article(article)
	if !c.isAppRequest(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.Article: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.Article: %w", err)
	}
	return nil
}

func (c Controller) Project(w http.ResponseWriter, r *http.Request) error {
	projectId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, strconv.IntSize)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			projectId = -1
		} else {
			return fmt.Errorf("controller.Project: %w", err)
		}
	}

	project, err := c.service.Project(int(projectId))
	if err != nil {
		return fmt.Errorf("controller.Project: %w", err)
	}

	pageComponent := page.Project(project)
	if !c.isAppRequest(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.Project: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.Project: %w", err)
	}
	return nil
}

func (c Controller) Serie(w http.ResponseWriter, r *http.Request) error {
	serieId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, strconv.IntSize)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			serieId = -1
		} else {
			return fmt.Errorf("controller.Serie: %w", err)
		}
	}

	// TODO: use workers
	serieContentParam := persistence.SerieContentQueryParam{Page: 1, Limit: 10}
	serie, err := c.service.Serie(int(serieId))
	if err != nil {
		return fmt.Errorf("controller<Controller.Serie>; %w", err)
	}
	serieProjects, err := c.service.SerieProjectList(serie.Id, serieContentParam)
	if err != nil {
		return fmt.Errorf("controller<Controller.Serie>; %w", err)
	}
	serieArticles, err := c.service.SerieArticleList(serie.Id, serieContentParam)
	if err != nil {
		return fmt.Errorf("controller<Controller.Serie>; %w", err)
	}

	pageComponent := page.Serie(serie, serieArticles, serieProjects)
	if !c.isAppRequest(r) {
		if err := c.serveWithBase(pageComponent, w, r); err != nil {
			return fmt.Errorf("controller.Serie: %w", err)
		}
	} else if err := pageComponent.Render(context.Background(), w); err != nil {
		return fmt.Errorf("controller.Serie: %w", err)
	}
	return nil
}
