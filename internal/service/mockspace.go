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

func (s Service) ArticleWritespace() (entity.Article, error) {
	f, err := os.Open("./static/testarticle.html")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return entity.Article{}, oops.NotFound{}
		}
		return entity.Article{}, fmt.Errorf("Service.ArticleWritespace: %w", err)
	}

	content, err := io.ReadAll(f)
	if err != nil {
		return entity.Article{}, fmt.Errorf("Service.ArticleWritespace: %w", err)
	}

	return entity.Article{
		Title:    "Testing article!",
		Subtitle: "A test article",
		Content:  string(content)}, nil
}

func (s Service) ProjectWritespace() (entity.Project, error) {
	f, err := os.Open("./static/testproject.html")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return entity.Project{}, oops.NotFound{}
		}
		return entity.Project{}, fmt.Errorf("Service.ProjectWritespace: %w", err)
	}

	content, err := io.ReadAll(f)
	if err != nil {
		return entity.Project{}, fmt.Errorf("Service.ProjectWritespace: %w", err)
	}

	return entity.Project{
		Name:        "Testing project!",
		Synopsis:    "A test project",
		Description: string(content)}, nil
}
