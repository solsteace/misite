package persistence

import (
	"fmt"
	"time"

	"github.com/solsteace/misite/internal/entity"
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
		return entity.Serie{}, fmt.Errorf(
			"persistence<Pg.Serie>: %w", err)
	}

	serie := entity.Serie{
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

func (p Pg) SerieArticleList(id int, param SerieContentQueryParam) ([]entity.SerieArticleList, error) {
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
		return []entity.SerieArticleList{}, fmt.Errorf(
			"persistence<Pg.SerieArticleList>: %w", err)
	}

	var serieArticles []entity.SerieArticleList
	for _, r := range rows {
		serieArticles = append(
			serieArticles, entity.SerieArticleList{
				Id:        r.Id,
				Title:     r.Title,
				Synopsis:  r.Synopsis,
				CreatedAt: r.CreatedAt,
				UpdatedAt: r.UpdatedAt,
			})
	}
	return serieArticles, nil
}

func (p Pg) SerieProjectList(id int, param SerieContentQueryParam) ([]entity.SerieProjectList, error) {
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
		return []entity.SerieProjectList{}, fmt.Errorf(
			"persistence<Pg.SerieProjectList>: %w", err)
	}

	var serieProjects []entity.SerieProjectList
	for _, r := range rows {
		serieProjects = append(
			serieProjects, entity.SerieProjectList{
				Id:        r.Id,
				Name:      r.Name,
				Synopsis:  r.Synopsis,
				CreatedAt: r.CreatedAt,
				UpdatedAt: r.UpdatedAt,
			})
	}
	return serieProjects, nil
}
