package persistence

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/solsteace/misite/internal/entity"
	"github.com/solsteace/misite/internal/utility/lib/oops"
)

type ArticlesQueryParam struct {
	Limit int
	Last  string
	Tag   []string
	Serie []string
}

func (p Pg) Articles(param ArticlesQueryParam) ([]entity.ArticleList, error) {
	query := `
		SELECT
			articles.id,
			articles.title,
			articles.subtitle,
			articles.created_at,
			articles.updated_at,
			tags.id AS "tag.id",
			tags.name AS "tag.name",
			series.id AS "serie.id",
			series.name AS "serie.name"
		FROM (
			SELECT *
			FROM articles
			WHERE
				id > $1 
				AND updated_at >= $2
				AND ($4::VARCHAR[] IS NULL
					OR EXISTS (
						SELECT 1
						FROM article_tags
						JOIN tags ON article_tags.tag_id = tags.id
						WHERE 
							article_tags.article_id = articles.id
							AND LOWER(tags.name) = ANY($4)
						GROUP BY article_id
						HAVING COUNT(DISTINCT tag_id) = CARDINALITY($4)))
				AND ($5::VARCHAR[] IS NULL 
					OR EXISTS(
						SELECT 1
						FROM articles AS temp_articles
						JOIN series ON series.id = articles.serie_id
						WHERE
							temp_articles.id = articles.id
							AND LOWER(series.name) = ANY($5)))
			ORDER BY
				updated_at DESC,
				id
			LIMIT $3
		) AS articles
		LEFT JOIN article_tags ON article_tags.article_id = articles.id
		LEFT JOIN tags ON article_tags.tag_id = tags.id
		LEFT JOIN series ON articles.serie_id = series.id
		ORDER BY 
			articles.updated_at DESC, 
			id`
	args := []any{
		0,               // $1 -> lastId
		time.Unix(0, 0), // $2 -> lastTime
		10,              // $3 -> limit
		nil,             // $4 -> tag filter
		nil,             // $5 -> serie filter
	}
	if tokens := strings.Split(param.Last, "-"); len(tokens) == 2 {
		if lastId, err := strconv.ParseInt(tokens[1], 10, strconv.IntSize); err == nil {
			args[0] = lastId
		}
		if lastTime, err := strconv.ParseInt(tokens[0], 10, strconv.IntSize); err == nil {
			args[1] = time.Unix(0, lastTime)
		}
	}
	if param.Limit > 0 {
		args[2] = param.Limit
	}
	if len(param.Tag) > 0 {
		args[3] = param.Tag
	}
	if len(param.Serie) > 0 {
		args[4] = param.Serie
	}

	var rows []struct {
		Id        int       `db:"id"`
		Title     string    `db:"title"`
		Subtitle  string    `db:"subtitle"`
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
		return []entity.ArticleList{}, fmt.Errorf(
			"persistence<Pg.Articles>: %s", err)
	} else if len(rows) == 0 {
		return []entity.ArticleList{}, nil
	}

	var articles []entity.ArticleList
	var lastArticle *entity.ArticleList
	var insertedTags map[int]struct{}
	for _, r := range rows {
		if lastArticle == nil || lastArticle.Id != r.Id {
			insertedTags = map[int]struct{}{}
			articles = append(articles, entity.ArticleList{
				Id:        r.Id,
				Title:     r.Title,
				Subtitle:  r.Subtitle,
				CreatedAt: r.CreatedAt,
				UpdatedAt: r.UpdatedAt})
			lastArticle = &articles[len(articles)-1]
		}
		if r.Serie.Id.Valid {
			lastArticle.Serie = &struct {
				Id   int
				Name string
			}{
				Id:   r.Serie.Id.V,
				Name: r.Serie.Name.V}
		}
		if r.Tag.Id.Valid {
			if _, ok := insertedTags[r.Tag.Id.V]; !ok {
				insertedTags[r.Tag.Id.V] = struct{}{}
				lastArticle.Tag = append(lastArticle.Tag, entity.Tag{
					Id:   r.Tag.Id.V,
					Name: r.Tag.Name.V})
			}
		}
	}
	return articles, nil
}

func (p Pg) Article(id int) (entity.Article, error) {
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
		return entity.Article{}, fmt.Errorf(
			"persistence<Pg.Article>: %w", err)
	} else if len(rows) == 0 {
		return entity.Article{}, fmt.Errorf(
			"persistence<Pg.Article>: %w", oops.NotFound{})
	}

	article := entity.Article{
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
