package persistence

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/solsteace/misite/internal/entity"
	"github.com/solsteace/misite/internal/utility/lib/oops"
)

func (p Pg) Article(id int) (entity.ArticlePage, error) {
	query := `
		SELECT
			articles.id AS "id",
			articles.title AS "title",
			articles.subtitle AS "subtitle",
			articles.content AS "content",
			articles.created_at AS "created_at",
			articles.updated_at AS "updated_at",
			tags.id AS "tag.id",
			tags.name AS "tag.name",
			series.id AS "serie.id",
			series.name AS "serie.name"
		FROM articles
		LEFT JOIN series ON articles.serie_id = series.id
		LEFT JOIN article_tags ON articles.id = article_tags.article_id
		LEFT JOIN tags ON article_tags.tag_id = tags.id
		WHERE articles.id = $1`
	args := []any{id}

	var rows []struct {
		Id        int       `db:"id"`
		Title     string    `db:"title"`
		Subtitle  string    `db:"subtitle"`
		Content   string    `db:"content"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`

		Serie struct {
			Id   sql.Null[int]    `db:"id"`
			Name sql.Null[string] `db:"name"`
		}
		Tag struct {
			Id   sql.Null[int]    `db:"id"`
			Name sql.Null[string] `db:"name"`
		}
	}
	if err := p.db.Select(&rows, query, args...); err != nil {
		return entity.ArticlePage{}, fmt.Errorf(
			"persistence<Pg.Article>: %w", err)
	} else if len(rows) == 0 {
		return entity.ArticlePage{}, fmt.Errorf(
			"persistence<Pg.Article>: %w", oops.NotFound{})
	}

	article := entity.ArticlePage{
		Id:        rows[0].Id,
		Title:     rows[0].Title,
		Subtitle:  rows[0].Subtitle,
		Content:   rows[0].Content,
		CreatedAt: rows[0].CreatedAt,
		UpdatedAt: rows[0].UpdatedAt}
	insertedTags := map[int]struct{}{}
	for _, r := range rows {
		if r.Serie.Id.Valid {
			article.Serie = &struct {
				Id   int
				Name string
			}{
				Id:   r.Serie.Id.V,
				Name: r.Serie.Name.V}
		}
		if r.Tag.Id.Valid {
			if _, ok := insertedTags[r.Tag.Id.V]; !ok {
				insertedTags[r.Tag.Id.V] = struct{}{}
				article.Tag = append(article.Tag, entity.Tag{
					Id:   r.Tag.Id.V,
					Name: r.Tag.Name.V})
			}
		}
	}
	return article, nil
}

func (p Pg) CountArticleMatchingTags(tagId []int) ([]entity.Tag, []int, error) {
	query := `
		WITH 
			tag_count_by_article AS (
				SELECT
					tag_id AS "id",
					COUNT(article_id) AS "count"
				FROM article_tags
				WHERE article_tags.tag_id = ANY($1::int[])
				GROUP BY tag_id)
		SELECT 
			tag_count.id,
			tag_count."count",
			tags.name AS "name"
		FROM tag_count_by_article AS tag_count
		JOIN tags ON tag_count.id = tags.id
		ORDER BY name`

	var rows []struct {
		Id      int    `db:"id"`
		Count   int    `db:"count"`
		TagName string `db:"name"`
	}
	args := []any{tagId}
	if err := p.db.Select(&rows, query, args...); err != nil {
		return []entity.Tag{}, []int{}, fmt.Errorf(
			"persistence<Pg.CountArticleMatchingTags>: %w", err)
	}

	var tags []entity.Tag
	var count []int
	for _, r := range rows {
		tags = append(tags, entity.Tag{
			Id:   r.Id,
			Name: r.TagName})
		count = append(count, r.Count)
	}
	return tags, count, nil
}

func (p Pg) Project(id int) (entity.ProjectPage, error) {
	query := `
		SELECT
			projects.id AS "id",
			projects.name AS "name",
			projects.synopsis AS "synopsis",
			projects.description AS "description",
			projects.created_at AS "created_at",
			projects.updated_at AS "updated_at",
			series.id AS "serie.id",
			series.name AS "serie.name",
			tags.id AS "tag.id",
			tags.name AS "tag.name",
			project_links.id AS "link.id",
			project_links.display_text AS "link.display_text",
			project_links.url AS "link.url"
		FROM projects
		LEFT JOIN project_links ON project_links.project_id = projects.id
		LEFT JOIN project_tags ON project_tags.project_id = projects.id
		LEFT JOIN tags ON tags.id = project_tags.tag_id
		LEFT JOIN series ON series.id = projects.devblog_serie
		WHERE projects.id = $1
		ORDER BY projects.id`
	args := []any{id}

	var rows []struct {
		Id          int       `db:"id"`
		Name        string    `db:"name"`
		Synopsis    string    `db:"synopsis"`
		Description string    `db:"description"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`

		Serie struct {
			Id   sql.Null[int]    `db:"id"`
			Name sql.Null[string] `db:"name"`
		}
		Tag struct {
			Id   sql.Null[int]    `db:"id"`
			Name sql.Null[string] `db:"name"`
		}
		Link struct {
			Id          sql.Null[int]    `db:"id"`
			DisplayText sql.Null[string] `db:"display_text"`
			Url         sql.Null[string] `db:"url"`
		}
	}
	if err := p.db.Select(&rows, query, args...); err != nil {
		return entity.ProjectPage{}, fmt.Errorf(
			"persistence<pg.Project>: %w", err)
	} else if len(rows) == 0 {
		return entity.ProjectPage{}, fmt.Errorf(
			"persistence<pg.Project>: %w", oops.NotFound{})
	}

	projectRow := rows[0]
	project := entity.ProjectPage{
		Id:          projectRow.Id,
		Name:        projectRow.Name,
		Synopsis:    projectRow.Synopsis,
		Description: projectRow.Description,
		CreatedAt:   projectRow.CreatedAt,
		UpdatedAt:   projectRow.UpdatedAt}
	if projectRow.Serie.Id.Valid {
		project.Serie = &struct {
			Id   int
			Name string
		}{
			Id:   projectRow.Serie.Id.V,
			Name: projectRow.Serie.Name.V}
	}

	insertedLink := map[int]struct{}{}
	insertedTag := map[int]struct{}{}
	for _, r := range rows {
		if r.Link.Id.Valid {
			if _, ok := insertedLink[r.Link.Id.V]; !ok {
				insertedLink[r.Link.Id.V] = struct{}{}
				project.Link = append(project.Link, struct {
					Id          int
					DisplayText string
					Url         string
				}{
					Id:          r.Link.Id.V,
					DisplayText: r.Link.DisplayText.V,
					Url:         r.Link.Url.V})
			}
		}
		if r.Tag.Id.Valid {
			if _, ok := insertedTag[r.Tag.Id.V]; !ok {
				insertedTag[r.Tag.Id.V] = struct{}{}
				project.Tag = append(project.Tag, struct {
					Id   int
					Name string
				}{
					Id:   r.Tag.Id.V,
					Name: r.Tag.Name.V})
			}
		}
	}
	return project, nil
}

func (p Pg) CountProjectMatchingTags(tagId []int) ([]entity.Tag, []int, error) {
	query := `
		WITH 
			tag_count_by_project AS (
				SELECT
					tag_id AS "id",
					COUNT(project_id) AS "count"
				FROM project_tags
				WHERE project_tags.tag_id = ANY($1::int[])
				GROUP BY tag_id)
		SELECT 
			tag_count.id,
			tag_count."count",
			tags.name AS "name"
		FROM tag_count_by_project AS tag_count
		JOIN tags ON tag_count.id = tags.id
		ORDER BY name`
	args := []any{tagId}

	var rows []struct {
		Id      int    `db:"id"`
		Count   int    `db:"count"`
		TagName string `db:"name"`
	}
	if err := p.db.Select(&rows, query, args...); err != nil {
		return []entity.Tag{}, []int{}, fmt.Errorf(
			"persistence<Pg.CountProjectMatchingTags>: %w", err)
	}

	var tags []entity.Tag
	var count []int
	for _, r := range rows {
		tags = append(tags, entity.Tag{
			Id:   r.Id,
			Name: r.TagName})
		count = append(count, r.Count)
	}
	return tags, count, nil
}

func (p Pg) Serie(id int) (entity.SeriePage, error) {
	var row struct {
		Id          int    `db:"id"`
		Name        string `db:"name"`
		Thumbnail   string `db:"thumbnail"`
		Description string `db:"description"`
		NArticles   int    `db:"n_articles"`
		NProjects   int    `db:"n_projects"`
	}
	query := `
		SELECT
			series.id,
			series.name,
			series.thumbnail,
			series.description,
			COUNT(DISTINCT articles.id) AS "n_articles",
			COUNT(DISTINCT projects.id) AS "n_projects"
		FROM series
		JOIN projects ON projects.devblog_serie = series.id
		JOIN articles ON articles.serie_id = series.id
		WHERE series.id = $1
		GROUP BY series.id`
	args := []any{id}
	if err := p.db.Get(&row, query, args...); err != nil {
		return entity.SeriePage{}, fmt.Errorf(
			"persistence<Pg.Serie>: %w", err)
	}

	serie := entity.SeriePage{
		Id:          row.Id,
		Name:        row.Name,
		Thumbnail:   row.Thumbnail,
		Description: row.Description,
		NArticle:    row.NArticles,
		NProject:    row.NProjects}
	return serie, nil
}

type SerieContentQueryParam struct {
	Page  int
	Limit int
}

func (p Pg) SerieArticleList(id int, param SerieContentQueryParam) ([]entity.SeriePageArticleList, error) {
	if param.Limit < 1 {
		param.Limit = 10
	}
	if param.Page < 1 {
		param.Page = 1
	}

	var rows []struct {
		Id        int       `db:"id"`
		Title     string    `db:"title"`
		Synopsis  string    `db:"synopsis"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	query := `
		SELECT
			id,
			title,
			subtitle AS "synopsis",
			created_at,
			updated_at
		FROM articles
		WHERE serie_id = $1
		ORDER BY serie_order DESC
		LIMIT $2 OFFSET $3`
	args := []any{
		id,
		param.Limit,
		(param.Page - 1) * param.Limit}
	if err := p.db.Select(&rows, query, args...); err != nil {
		return []entity.SeriePageArticleList{}, fmt.Errorf(
			"persistence<Pg.SerieArticleList>: %w", err)
	}

	var serieArticles []entity.SeriePageArticleList
	for _, r := range rows {
		serieArticles = append(
			serieArticles, entity.SeriePageArticleList{
				Id:        r.Id,
				Title:     r.Title,
				Synopsis:  r.Synopsis,
				CreatedAt: r.CreatedAt,
				UpdatedAt: r.UpdatedAt,
			})
	}
	return serieArticles, nil
}

func (p Pg) SerieProjectList(id int, param SerieContentQueryParam) ([]entity.SeriePageProjectList, error) {
	if param.Limit < 1 {
		param.Limit = 10
	}
	if param.Page < 1 {
		param.Page = 1
	}
	var rows []struct {
		Id        int       `db:"id"`
		Name      string    `db:"name"`
		Synopsis  string    `db:"synopsis"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	query := `
		SELECT
			id,
			name,
			synopsis,
			created_at,
			updated_at
		FROM projects
		WHERE devblog_serie = $1
		LIMIT $2 OFFSET $3`
	args := []any{
		id,
		param.Limit,
		(param.Page - 1) * param.Limit}
	if err := p.db.Select(&rows, query, args...); err != nil {
		return []entity.SeriePageProjectList{}, fmt.Errorf(
			"persistence<Pg.SerieProjectList>: %w", err)
	}

	var serieProjects []entity.SeriePageProjectList
	for _, r := range rows {
		serieProjects = append(
			serieProjects, entity.SeriePageProjectList{
				Id:        r.Id,
				Name:      r.Name,
				Synopsis:  r.Synopsis,
				CreatedAt: r.CreatedAt,
				UpdatedAt: r.UpdatedAt,
			})
	}
	return serieProjects, nil
}
