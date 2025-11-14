package service

import (
	"fmt"

	"github.com/solsteace/misite/internal/entity"
)

func (s Service) Project(id int) (entity.Project, error) {
	project, err := s.store.Project(id)
	if err != nil {
		return entity.Project{}, fmt.Errorf(
			"service<Service.Project>: %w", err)
	}
	return project, nil
}

func (s Service) Article(id int) (entity.Article, error) {
	article, err := s.store.Article(id)
	if err != nil {
		return entity.Article{}, fmt.Errorf(
			"service<Service.Article>: %w", err)
	}
	return article, nil
}

func (s Service) Serie(id int) (entity.Serie, error) {
	serie, err := s.store.Serie(id)
	if err != nil {
		return entity.Serie{}, fmt.Errorf(
			"service<Service.Serie>: %w", err)
	}
	return serie, nil
}
