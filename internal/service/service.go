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

func (s Service) Articles(page, limit int) ([]entity.Article, error) {
	articles, err := s.store.Articles(persistence.ArticlesQueryParam{
		Page: page, Limit: limit})
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

func (s Service) Projects(page, limit int) ([]entity.Project, error) {
	projects, err := s.store.Projects(persistence.ProjectsQueryParam{
		Page: page, Limit: limit})
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
