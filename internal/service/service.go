package service

import (
	"fmt"

	"github.com/solsteace/misite/internal/entity"
	"github.com/solsteace/misite/internal/persistence"
)

type Service struct {
	store *persistence.Pg // TODO: change to interface if needed
}

func NewService(store *persistence.Pg) Service {
	return Service{store: store}
}

func (s Service) Articles(param persistence.ArticlesQueryParam) ([]entity.Article, error) {
	articles, err := s.store.Articles(param)
	if err != nil {
		return []entity.Article{}, fmt.Errorf(
			"service<Service.Articles>: %w", err)
	}
	return articles, nil
}

func (s Service) Article(id int) (entity.Article, error) {
	article, err := s.store.Article(id)
	if err != nil {
		return entity.Article{}, fmt.Errorf(
			"service<Service.Article>: %w", err)
	}
	return article, nil
}

func (s Service) Projects(param persistence.ProjectsQueryParam) ([]entity.Project, error) {
	projects, err := s.store.Projects(param)
	if err != nil {
		return []entity.Project{}, fmt.Errorf(
			"service<Service.Projects>: %w", err)
	}
	return projects, nil
}

func (s Service) Project(id int) (entity.Project, error) {
	project, err := s.store.Project(id)
	if err != nil {
		return entity.Project{}, fmt.Errorf(
			"service<Service.Project>: %w", err)
	}
	return project, nil
}
