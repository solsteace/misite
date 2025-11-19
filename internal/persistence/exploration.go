package persistence

import (
	"fmt"

	"github.com/solsteace/misite/internal/entity"
)

type TagQueryParams struct {
	Page  int
	Limit int
}

func (p Pg) ArticleTags(param TagQueryParams) ([]entity.TagStat, error) {
	query := `
		SELECT 
			tags.id,
			tags.name,
			COUNT(article_tags.article_id) AS "count"
		FROM tags
		JOIN article_tags ON article_tags.tag_id = tags.id
		GROUP BY tags.id
		ORDER BY tags.name
		LIMIT $2 OFFSET $1`
	args := []any{
		0,  // $1 -> offset
		10} // $2 -> limit
	if param.Page > 0 {
		args[0] = (param.Page - 1) * param.Limit
	}
	if param.Limit > 0 {
		args[1] = param.Limit
	}

	var rows []struct {
		Id    int    `db:"id"`
		Name  string `db:"name"`
		Count int    `db:"count"`
	}
	if err := p.db.Select(&rows, query, args...); err != nil {
		return []entity.TagStat{}, fmt.Errorf(
			"persistence<Pg.ArticleTags>: %w", err)
	}

	var tagStat []entity.TagStat
	for _, r := range rows {
		tagStat = append(tagStat, entity.TagStat{
			Id:    r.Id,
			Name:  r.Name,
			Count: r.Count})
	}
	return tagStat, nil
}

func (p Pg) ProjectTags(param TagQueryParams) ([]entity.TagStat, error) {
	query := `
		SELECT
			tags.id,
			tags.name,
			COUNT(project_tags.project_id) AS "count"
		FROM tags
		JOIN project_tags ON project_tags.tag_id = tags.id
		GROUP BY tags.id
		ORDER BY tags.name
		LIMIT $2 OFFSET $1`
	args := []any{
		0,  // $1 -> offset
		10} // $2 -> limit
	if param.Page > 0 {
		args[0] = (param.Page - 1) * param.Limit
	}
	if param.Limit > 0 {
		args[1] = param.Limit
	}

	var rows []struct {
		Id    int    `db:"id"`
		Name  string `db:"name"`
		Count int    `db:"count"`
	}
	if err := p.db.Select(&rows, query, args...); err != nil {
		return []entity.TagStat{}, fmt.Errorf(
			"persistence<Pg.ProjectTags>: %w", err)
	}

	var tagStat []entity.TagStat
	for _, r := range rows {
		tagStat = append(tagStat, entity.TagStat{
			Id:    r.Id,
			Name:  r.Name,
			Count: r.Count})
	}
	return tagStat, nil
}
