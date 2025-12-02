package service

import (
	"fmt"

	"github.com/solsteace/misite/internal/entity"
	"github.com/solsteace/misite/internal/persistence"
	"github.com/solsteace/misite/internal/utility/lib/oops"
)

func (s Service) Projects(param persistence.ProjectsQueryParam) ([]entity.ProjectListPage, error) {
	projects, err := s.store.Projects(param)
	if err != nil {
		return []entity.ProjectListPage{}, fmt.Errorf(
			"service<Service.Projects>: %w", err)
	}
	return projects, nil
}

func (s Service) Articles(param persistence.ArticlesQueryParam) ([]entity.ArticleListPage, error) {
	articles, err := s.store.Articles(param)
	if err != nil {
		return []entity.ArticleListPage{}, fmt.Errorf(
			"service<Service.Articles>: %w", err)
	}
	return articles, nil
}

func (s Service) SerieList(param persistence.SerieListQueryParam) ([]entity.SerieListPage, error) {
	serieList, err := s.store.SerieList(param)
	if err != nil {
		return []entity.SerieListPage{}, fmt.Errorf(
			"service<Service.SerieList>: %w", err)
	}
	return serieList, nil
}

func (s Service) SerieArticleList(
	id int,
	param persistence.SerieContentQueryParam,
) ([]entity.SeriePageArticleList, error) {
	serieArticles, err := s.store.SerieArticleList(id, param)
	if err != nil {
		return []entity.SeriePageArticleList{}, fmt.Errorf(
			"service<Service.SerieArticles>: %w", err)
	}
	return serieArticles, nil
}

func (s Service) SerieProjectList(
	id int,
	param persistence.SerieContentQueryParam,
) ([]entity.SeriePageProjectList, error) {
	serieProjects, err := s.store.SerieProjectList(id, param)
	if err != nil {
		return []entity.SeriePageProjectList{}, fmt.Errorf(
			"service<Service.SerieProjects>: %w", err)
	}
	return serieProjects, nil
}

func (s Service) Tags(by string, param persistence.TagQueryParams) ([]entity.TagStatPage, error) {
	var tagStats []entity.TagStatPage
	var err error

	switch by {
	case "article":
		tagStats, err = s.store.ArticleTags(param)
	case "project":
		tagStats, err = s.store.ProjectTags(param)
	default:
		return []entity.TagStatPage{}, oops.NotFound{
			Err: fmt.Errorf("`by` should be either `article` or `project`")}
	}
	if err != nil {
		return []entity.TagStatPage{}, fmt.Errorf(
			"service<Service.Tags>: %w", err)
	}

	return tagStats, nil
}
