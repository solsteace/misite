package persistence

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/solsteace/misite/internal/entity"
	"github.com/solsteace/misite/internal/utility/lib/oops"
)

type SerieListQueryParam struct {
	Page  int
	Limit int
}

func (p Pg) SerieList(param SerieListQueryParam) ([]entity.SerieList, error) {
	if param.Limit < 1 {
		param.Limit = 10
	}
	if param.Page < 1 {
		param.Page = 1
	}

	var rows []struct {
		Id          int       `db:"id"`
		Name        string    `db:"name"`
		Description string    `db:"description"`
		CreatedAt   time.Time `db:"created_at"`
	}
	query := `
		SELECT
			series.id,
			series.name,
			series.description,
			series.created_at
		FROM series
		ORDER BY series.created_at
		LIMIT $1 OFFSET $2`
	args := []any{
		param.Limit,
		(param.Page - 1) * param.Limit}
	if err := p.db.Select(&rows, query, args...); err != nil {
		return []entity.SerieList{}, fmt.Errorf(
			"persistence<Pg.Series>: %w", err)
	}

	var serieList []entity.SerieList
	var last *entity.SerieList
	for _, r := range rows {
		if last == nil || last.Id != r.Id {
			sl := entity.SerieList{
				Id:          r.Id,
				Name:        r.Name,
				Description: r.Description,
				CreatedAt:   r.CreatedAt}
			serieList = append(serieList, sl)
			last = &serieList[len(serieList)-1]
		}
	}
	return serieList, nil
}

func (p Pg) Serie(id int) (entity.Serie, error) {
	var rows []struct {
		Id          int    `db:"id"`
		Name        string `db:"name"`
		Thumbnail   string `db:"thumbnail"`
		Description string `db:"description"`

		Article struct {
			Id        sql.Null[int]       `db:"id"`
			Title     sql.Null[string]    `db:"title"`
			Synopsis  sql.Null[string]    `db:"synopsis"`
			CreatedAt sql.Null[time.Time] `db:"created_at"`
			UpdatedAt sql.Null[time.Time] `db:"updated_at"`
		}
		Project struct {
			Id        sql.Null[int]       `db:"id"`
			Name      sql.Null[string]    `db:"name"`
			Synopsis  sql.Null[string]    `db:"synopsis"`
			CreatedAt sql.Null[time.Time] `db:"created_at"`
			UpdatedAt sql.Null[time.Time] `db:"updated_at"`
		}
	}
	query := `
		SELECT 	
			series.id,
			series.name,
			series.thumbnail,
			series.description,
			articles.id AS "article.id",
			articles.title AS "article.title",
			articles.subtitle AS "article.synopsis",
			articles.created_at AS "article.created_at",
			articles.updated_at AS "article.updated_at",
			articles.created_at AS "article.created_at",
			projects.id AS "project.id",
			projects.name AS "project.name",
			projects.synopsis AS "project.synopsis",
			projects.updated_at AS "project.updated_at",
			projects.created_at AS "project.created_at"
		FROM series
		LEFT JOIN projects ON projects.devblog_serie = series.id
		LEFT JOIN articles ON articles.serie_id = series.id
		WHERE series.id = $1
		ORDER BY articles.serie_order`
	args := []any{id}
	if err := p.db.Select(&rows, query, args...); err != nil {
		return entity.Serie{}, fmt.Errorf(
			"persistence<Pg.Serie>: %w", err)
	} else if len(rows) == 0 {
		return entity.Serie{}, oops.NotFound{}
	}

	insertedArticles := map[int]struct{}{}
	insertedProjects := map[int]struct{}{}
	serie := entity.Serie{
		Id:          rows[0].Id,
		Name:        rows[0].Name,
		Thumbnail:   rows[0].Thumbnail,
		Description: rows[0].Description}
	for _, r := range rows {
		if r.Article.Id.Valid {
			article := r.Article
			if _, ok := insertedArticles[article.Id.V]; !ok {
				insertedArticles[article.Id.V] = struct{}{}
				serie.Article = append(serie.Article,
					struct {
						Id        int
						Title     string
						Synopsis  string
						CreatedAt time.Time
						UpdatedAt time.Time
					}{
						article.Id.V,
						article.Title.V,
						article.Synopsis.V,
						article.CreatedAt.V,
						article.UpdatedAt.V})
			}
		}
		if r.Project.Id.Valid {
			project := r.Project
			if _, ok := insertedProjects[project.Id.V]; !ok {
				insertedProjects[project.Id.V] = struct{}{}
				serie.Project = append(serie.Project, struct {
					Id        int
					Name      string
					Synopsis  string
					CreatedAt time.Time
					UpdatedAt time.Time
				}{
					project.Id.V,
					project.Name.V,
					project.Synopsis.V,
					project.CreatedAt.V,
					project.UpdatedAt.V})
			}
		}
	}
	return serie, nil
}
