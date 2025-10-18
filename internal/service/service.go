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
		return []entity.Article{}, err
	}
	return articles, nil
}

func (s Service) Article(id int) (entity.Article, error) {
	article, err := s.store.Article(id)
	if err != nil {
		return entity.Article{}, err
	}
	fmt.Printf("%+v", article)
	return article, nil
}
