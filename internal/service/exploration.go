package service

import (
	"fmt"

	"github.com/solsteace/misite/internal/entity"
	"github.com/solsteace/misite/internal/persistence"
)

func (s Service) Projects(param persistence.ProjectsQueryParam) ([]entity.ProjectList, error) {
	projects, err := s.store.Projects(param)
	if err != nil {
		return []entity.ProjectList{}, fmt.Errorf(
			"service<Service.Projects>: %w", err)
	}
	return projects, nil
}

func (s Service) Articles(param persistence.ArticlesQueryParam) ([]entity.ArticleList, error) {
	articles, err := s.store.Articles(param)
	if err != nil {
		return []entity.ArticleList{}, fmt.Errorf(
			"service<Service.Articles>: %w", err)
	}
	return articles, nil
}

func (s Service) SerieList(param persistence.SerieListQueryParam) ([]entity.SerieList, error) {
	serieList, err := s.store.SerieList(param)
	if err != nil {
		return []entity.SerieList{}, fmt.Errorf(
			"service<Service.SerieList>: %w", err)
	}
	return serieList, nil
}

func (s Service) SerieArticleList(
	id int,
	param persistence.SerieContentQueryParam,
) ([]entity.SerieArticleList, error) {
	serieArticles, err := s.store.SerieArticleList(id, param)
	if err != nil {
		return []entity.SerieArticleList{}, fmt.Errorf(
			"service<Service.SerieArticles>: %w", err)
	}
	return serieArticles, nil
}

func (s Service) SerieProjectList(
	id int,
	param persistence.SerieContentQueryParam,
) ([]entity.SerieProjectList, error) {
	serieProjects, err := s.store.SerieProjectList(id, param)
	if err != nil {
		return []entity.SerieProjectList{}, fmt.Errorf(
			"service<Service.SerieProjects>: %w", err)
	}
	return serieProjects, nil
}
