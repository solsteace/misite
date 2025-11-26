package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/solsteace/misite/internal/component/page"
	"github.com/solsteace/misite/internal/utility/lib/oops"
)

func (c Controller) MockSpace(w http.ResponseWriter, r *http.Request) error {
	urlQuery := r.URL.Query()
	switch urlQuery.Get("for") {
	case "article":
		article, err := c.service.ArticleWritespace()
		if err != nil {
			return fmt.Errorf("controller.Writespace: %w", err)
		}

		pageComponent := page.Article(article)
		if !c.isAppRequest(r) {
			if err := c.serveWithBase(pageComponent, w, r); err != nil {
				return fmt.Errorf("controller.Writespace: %w", err)
			}
		} else if err := pageComponent.Render(context.Background(), w); err != nil {
			return fmt.Errorf("controller.Writespace: %w", err)
		}
		return nil
	case "project":
		project, err := c.service.ProjectWritespace()
		if err != nil {
			return fmt.Errorf("controller.Writespace: %w", err)
		}
		pageComponent := page.Project(project)
		if !c.isAppRequest(r) {
			if err := c.serveWithBase(pageComponent, w, r); err != nil {
				return fmt.Errorf("controller.Writespace: %w", err)
			}
		} else if err := pageComponent.Render(context.Background(), w); err != nil {
			return fmt.Errorf("controller.Writespace: %w", err)
		}
		return nil

	}
	return oops.NotFound{}
}
