package service

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/solsteace/misite/internal/entity"
	"github.com/solsteace/misite/internal/persistence"
)

type Service struct {
	store *persistence.Pg // TODO: change to interface if needed
}

func NewService(store *persistence.Pg) Service {
	return Service{store: store}
}

func (s Service) MostTagsOnArticles(articles []entity.ArticleList, limit int) []int {
	tags := make(map[int]int)
	tagIds := []int{}
	for _, a := range articles {
		for _, t := range a.Tag {
			if _, ok := tags[t.Id]; !ok {
				tags[t.Id] = 0
				tagIds = append(tagIds, t.Id)
			}
			tags[t.Id] += 1
		}
	}

	slices.SortFunc(tagIds, func(a, b int) int {
		tagCount1 := tags[a]
		tagCount2 := tags[b]
		if res := cmp.Compare(tagCount1, tagCount2); res != 0 {
			return -res
		}
		return -1
	})

	if len(tagIds) > limit {
		return tagIds[:limit]
	}
	return tagIds
}

func (s Service) CountArticleMatchingTags(tagIds []int) ([]entity.Tag, []int, error) {
	tags, count, err := s.store.CountArticleMatchingTags(tagIds)
	if err != nil {
		return []entity.Tag{}, []int{}, fmt.Errorf(
			"service<Service.CountArticleMatchingTags>: %w", err)
	}
	return tags, count, nil
}

func (s Service) MostTagsOnProjects(projects []entity.ProjectList, limit int) []int {
	tags := make(map[int]int)
	tagIds := []int{}
	for _, p := range projects {
		for _, t := range p.Tag {
			if _, ok := tags[t.Id]; !ok {
				tags[t.Id] = 0
				tagIds = append(tagIds, t.Id)
			}
			tags[t.Id] += 1
		}
	}

	slices.SortFunc(tagIds, func(a, b int) int {
		tagCount1 := tags[a]
		tagCount2 := tags[b]
		if res := cmp.Compare(tagCount1, tagCount2); res != 0 {
			return -res
		}
		return -1
	})

	if len(tagIds) > limit {
		return tagIds[:limit]
	}
	return tagIds
}

func (s Service) CountProjectMatchingTags(tagIds []int) ([]entity.Tag, []int, error) {
	tags, count, err := s.store.CountProjectMatchingTags(tagIds)
	if err != nil {
		return []entity.Tag{}, []int{}, fmt.Errorf(
			"service<Service.CountProjectMatchingTags>: %w", err)
	}
	return tags, count, nil
}
