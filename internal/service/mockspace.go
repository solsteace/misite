package service

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/solsteace/misite/internal/entity"
	"github.com/solsteace/misite/internal/utility/lib/oops"
)

func (s Service) ArticleWritespace() (entity.ArticlePage, error) {
	f, err := os.Open("./static/testarticle.html")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return entity.ArticlePage{}, oops.NotFound{}
		}
		return entity.ArticlePage{}, fmt.Errorf("Service.ArticleWritespace: %w", err)
	}

	content, err := io.ReadAll(f)
	if err != nil {
		return entity.ArticlePage{}, fmt.Errorf("Service.ArticleWritespace: %w", err)
	}

	return entity.ArticlePage{
		Title:    "Testing article!",
		Subtitle: "A test article",
		Content:  string(content)}, nil
}

func (s Service) ProjectWritespace() (entity.ProjectPage, error) {
	f, err := os.Open("./static/testproject.html")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return entity.ProjectPage{}, oops.NotFound{}
		}
		return entity.ProjectPage{}, fmt.Errorf("Service.ProjectWritespace: %w", err)
	}

	content, err := io.ReadAll(f)
	if err != nil {
		return entity.ProjectPage{}, fmt.Errorf("Service.ProjectWritespace: %w", err)
	}

	return entity.ProjectPage{
		Name:        "Testing project!",
		Synopsis:    "A test project",
		Description: string(content)}, nil
}
