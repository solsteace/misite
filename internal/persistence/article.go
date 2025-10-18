package persistence

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
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
		LEFT JOIN article_series ON articles.id = article_series.article_id
		LEFT JOIN series ON article_series.serie_id = series.id
		LEFT JOIN article_tags ON articles.id = article_tags.article_id
		LEFT JOIN tags ON article_tags.tag_id = tags.id
		ORDER BY articles.id
		LIMIT $1
		OFFSET $2`
	args := []any{
		param.Limit,
		(param.Page - 1) * param.Limit}
	fmt.Println(args...)

	if len(param.TagId) > 0 {
		query = fmt.Sprintf("%s WHERE article_tags.id IN (?)", query)
		queryTagFilter, argsTagFilter, err := sqlx.In(query, param.TagId)
		if err != nil {
			return []entity.Article{}, err
		}

		query = queryTagFilter
		args = append(args, argsTagFilter...)
	}

	if len(param.SerieId) > 0 {
		query = fmt.Sprintf("%s WHERE article_series.id IN (?)", query)
		querySerieFilter, argsSerieFilter, err := sqlx.In(query, param.SerieId)
		if err != nil {
			return []entity.Article{}, err
		}

		query = querySerieFilter
		args = append(args, argsSerieFilter...)
	}

	var rows []pgArticles
	if err := p.db.Select(&rows, p.db.Rebind(query), args...); err != nil {
		return []entity.Article{}, err
	}
	if len(rows) == 0 {
		return []entity.Article{}, nil
	}

	var articles []entity.Article
	var article entity.Article
	for i, r := range rows {
		if article.Id != r.Id {
			if i != 0 {
				articles = append(articles, article)
			}

			article = entity.Article{
				Id:        r.Id,
				Title:     r.Title,
				Subtitle:  r.Subtitle,
				CreatedAt: r.CreatedAt}

			if r.Serie.Id.Valid {
				article.Serie = &entity.Serie{
					Id:   r.Serie.Id.V,
					Name: r.Serie.Name.V}
			}
		}

		if r.Tag.Id.Valid {
			article.Tag = append(article.Tag, entity.Tag{
				Id:   r.Tag.Id.V,
				Name: r.Tag.Name.V})
		}
	}
	articles = append(articles, article)
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
	rows := []pgArticle{}
	args := []any{id}
	if err := p.db.Select(&rows, query, args...); err != nil {
		return entity.Article{}, err
	}
	if len(rows) == 0 {
		return entity.Article{}, oops.NotFound{}
	}

	articleRow := rows[0]
	article := entity.Article{
		Id:        articleRow.Id,
		Title:     articleRow.Title,
		Subtitle:  articleRow.Subtitle,
		Content:   articleRow.Content,
		CreatedAt: articleRow.CreatedAt}
	for _, r := range rows {
		if r.Serie.Valid {
			article.Serie = &entity.Serie{
				Id: r.Serie.V}
		}

		if r.Tag.Id.Valid {
			article.Tag = append(article.Tag, entity.Tag{
				Id:   r.Tag.Id.V,
				Name: r.Tag.Name.V})
		}
	}
	return article, nil
}
