package persistence

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/solsteace/misite/internal/entity"
	"github.com/solsteace/misite/internal/utility/lib/oops"
)

type ArticlesQueryParam struct {
	Page    int
	Limit   int
	TagId   []int
	SerieId []int
}

type pgArticles struct {
	Id        int       `db:"id"`
	Title     string    `db:"title"`
	Subtitle  string    `db:"subtitle"`
	Thumbnail string    `db:"thumbnail"`
	CreatedAt time.Time `db:"created_at"`

	Serie struct {
		Id   sql.Null[int]    `db:"id"`
		Name sql.Null[string] `db:"name"`
	}
	Tag struct {
		Id   sql.Null[int]    `db:"id"`
		Name sql.Null[string] `db:"name"`
	}
}

func (p Pg) Articles(param ArticlesQueryParam) ([]entity.Article, error) {
	if param.Limit < 1 {
		param.Limit = 10
	}
	if param.Page < 1 {
		param.Page = 1
	}

	query := `
		WITH 
			matching_articles_by_tag AS ( 
				SELECT 
					article_id, 
					COUNT(DISTINCT tag_id) AS "n_tag"
				FROM article_tags
				WHERE 
					($1::int[] IS NULL) 
					OR tag_id = ANY($1)
				GROUP BY article_id
				HAVING
					($1::int[] IS NULL)
					OR COUNT(tag_id) = CARDINALITY($1)
				LIMIT $3 OFFSET $4)
		SELECT
			articles.id AS "id",
			articles.title AS "title",
			articles.subtitle AS "subtitle",
			articles.thumbnail AS "thumbnail",
			articles.created_at AS "created_at",
			tags.id AS "tag.id",
			tags.name AS "tag.name",
			series.id AS "serie.id",
			series.name AS "serie.name"
		FROM articles
		LEFT JOIN article_tags ON article_tags.article_id = articles.id
		LEFT JOIN tags ON article_tags.tag_id = tags.id
		LEFT JOIN article_series ON article_series.article_id = articles.id
		LEFT JOIN series ON article_series.serie_id = series.id
		WHERE 
			(SELECT true
				FROM matching_articles_by_tag
				WHERE matching_articles_by_tag.article_id = articles.id)
			AND ($2::int[] IS NULL OR article_series.serie_id = ANY($2))
		ORDER BY articles.id`
	args := []any{
		nil,
		nil,
		param.Limit,
		(param.Page - 1) * param.Limit}
	if len(param.TagId) > 0 {
		args[0] = param.TagId
	}
	if len(param.SerieId) > 0 {
		args[1] = param.SerieId
	}

	var rows []pgArticles
	if err := p.db.Select(&rows, query, args...); err != nil {
		return []entity.Article{}, fmt.Errorf(
			"persistence<Pg.Articles>: %s", err)
	} else if len(rows) == 0 {
		return []entity.Article{}, nil
	}

	var articles []entity.Article
	var lastArticle *entity.Article
	var insertedTags map[int]struct{}
	for _, r := range rows {
		if lastArticle == nil || lastArticle.Id != r.Id {
			insertedTags = map[int]struct{}{}
			articles = append(articles, entity.Article{
				Id:        r.Id,
				Title:     r.Title,
				Subtitle:  r.Subtitle,
				Thumbnail: r.Thumbnail,
				CreatedAt: r.CreatedAt})
			lastArticle = &articles[len(articles)-1]
		}
		if r.Serie.Id.Valid {
			lastArticle.Serie = &entity.Serie{
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

type pgArticle struct {
	Id        int       `db:"id"`
	Title     string    `db:"title"`
	Subtitle  string    `db:"subtitle"`
	Content   string    `db:"content"`
	Thumbnail string    `db:"thumbnail"`
	CreatedAt time.Time `db:"created_at"`

	Serie sql.Null[int] `db:"serie.id"`
	Tag   struct {
		Id   sql.Null[int]    `db:"id"`
		Name sql.Null[string] `db:"name"`
	}
}

func (p Pg) Article(id int) (entity.Article, error) {
	query := `
		SELECT
			articles.id AS "id",
			articles.title AS "title",
			articles.subtitle AS "subtitle",
			articles.content AS "content",
			articles.thumbnail AS "thumbnail",
			articles.created_at AS "created_at",
			tags.id AS "tag.id",
			tags.name AS "tag.name",
			series.id AS "serie.id"
		FROM articles
		LEFT JOIN article_series ON articles.id = article_series.article_id
		LEFT JOIN series ON article_series.serie_id = series.id
		LEFT JOIN article_tags ON articles.id = article_tags.article_id
		LEFT JOIN tags ON article_tags.tag_id = tags.id
		WHERE articles.id = $1
		ORDER BY articles.id`
	args := []any{id}
	var rows []pgArticle
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
		Thumbnail: rows[0].Thumbnail,
		CreatedAt: rows[0].CreatedAt}
	insertedTags := map[int]struct{}{}
	for _, r := range rows {
		if r.Serie.Valid {
			article.Serie = &entity.Serie{
				Id: r.Serie.V}
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
