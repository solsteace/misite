package persistence

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/solsteace/misite/internal/entity"
)

type ArticlesQueryParam struct {
	Limit int
	Last  string
	Tag   []string
	Serie []string
}

func (p Pg) Articles(param ArticlesQueryParam) ([]entity.ArticleListPage, error) {
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
		return []entity.ArticleListPage{}, fmt.Errorf(
			"persistence<Pg.Articles>: %s", err)
	} else if len(rows) == 0 {
		return []entity.ArticleListPage{}, nil
	}

	var articles []entity.ArticleListPage
	var lastArticle *entity.ArticleListPage
	var insertedTags map[int]struct{}
	for _, r := range rows {
		if lastArticle == nil || lastArticle.Id != r.Id {
			insertedTags = map[int]struct{}{}
			articles = append(articles, entity.ArticleListPage{
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

type ProjectsQueryParam struct {
	Limit int
	Last  string
	Tag   []string
	Serie []string
}

func (p Pg) Projects(param ProjectsQueryParam) ([]entity.ProjectListPage, error) {
	query := `
		SELECT
		   	projects.id,
		   	projects.name,
		   	projects.synopsis,
		   	projects.created_at,
		   	projects.updated_at,
		   	tags.id AS "tag.id",
		   	tags.name AS "tag.name",
		   	series.id AS "serie.id",
		   	series.name AS "serie.name"
		FROM (
			SELECT *
			FROM projects
			WHERE
				id > $1
				AND updated_at >= $2
				AND ($4::VARCHAR[] IS NULL
					OR EXISTS (
						SELECT 1
						FROM project_tags
						JOIN tags ON project_tags.tag_id = tags.id
						WHERE
							project_tags.project_id = projects.id
							AND LOWER(tags.name) = ANY($4)
						GROUP BY project_id
						HAVING COUNT(DISTINCT tag_id) = CARDINALITY($4)))
				AND ($5::VARCHAR[] IS NULL
					OR EXISTS (
						SELECT 1
						FROM projects AS temp_projects
						JOIN series ON series.id = projects.devblog_serie
						WHERE 
							temp_projects.id = projects.id
							AND LOWER(series.name) = ANY($5)))
			ORDER BY
				updated_at DESC,
				id
			LIMIT $3
		) AS projects
		LEFT JOIN project_tags ON project_tags.project_id = projects.id
		LEFT JOIN tags ON project_tags.tag_id = tags.id
		LEFT JOIN series ON projects.devblog_serie = series.id
		ORDER BY 
			projects.updated_at DESC, 
			id`
	args := []any{
		0,               // $1 -> lastId
		time.Unix(0, 0), // $2 -> lastTime
		10,              // $3 -> limit
		nil,             // $4 -> tagList
		nil}             // $5 -> serieList
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
		Name      string    `db:"name"`
		Thumbnail string    `db:"thumbnail"`
		Synopsis  string    `db:"synopsis"`
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
		return []entity.ProjectListPage{}, fmt.Errorf(
			"persistence<Pg.Projects>: %w", err)
	} else if len(rows) == 0 {
		return []entity.ProjectListPage{}, nil
	}

	var projects []entity.ProjectListPage
	var lastProject *entity.ProjectListPage
	var insertedTag map[int]struct{}
	for _, r := range rows {
		if lastProject == nil || lastProject.Id != r.Id {
			insertedTag = map[int]struct{}{}
			projects = append(projects, entity.ProjectListPage{
				Id:        r.Id,
				Name:      r.Name,
				Synopsis:  r.Synopsis,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now()})
			lastProject = &projects[len(projects)-1]
		}
		if r.Tag.Id.Valid {
			if _, ok := insertedTag[r.Tag.Id.V]; !ok {
				insertedTag[r.Tag.Id.V] = struct{}{}
				lastProject.Tag = append(lastProject.Tag, struct {
					Id   int
					Name string
				}{
					Id:   r.Tag.Id.V,
					Name: r.Tag.Name.V})
			}
		}
		if r.Serie.Id.Valid && lastProject.Serie == nil {
			lastProject.Serie = &struct {
				Id   int
				Name string
			}{
				r.Serie.Id.V,
				r.Serie.Name.V}
		}
	}
	return projects, nil
}

type TagQueryParams struct {
	Page  int
	Limit int
}

func (p Pg) ArticleTags(param TagQueryParams) ([]entity.TagStatPage, error) {
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
		return []entity.TagStatPage{}, fmt.Errorf(
			"persistence<Pg.ArticleTags>: %w", err)
	}

	var tagStat []entity.TagStatPage
	for _, r := range rows {
		tagStat = append(tagStat, entity.TagStatPage{
			Id:    r.Id,
			Name:  r.Name,
			Count: r.Count})
	}
	return tagStat, nil
}

func (p Pg) ProjectTags(param TagQueryParams) ([]entity.TagStatPage, error) {
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
		return []entity.TagStatPage{}, fmt.Errorf(
			"persistence<Pg.ProjectTags>: %w", err)
	}

	var tagStat []entity.TagStatPage
	for _, r := range rows {
		tagStat = append(tagStat, entity.TagStatPage{
			Id:    r.Id,
			Name:  r.Name,
			Count: r.Count})
	}
	return tagStat, nil
}

type SerieListQueryParam struct {
	Title string
	Last  string
	Limit int
}

func (p Pg) SerieList(param SerieListQueryParam) ([]entity.SerieListPage, error) {
	query := `
		SELECT
			id,
			name,
			description,
			created_at
		FROM series
		WHERE 
			id > $1 
			AND created_at >= $2
			AND LOWER(name) LIKE $4
		ORDER BY 
			created_at DESC,
			id
		LIMIT $3`
	args := []any{
		0,               // $1 -> lastId
		time.Unix(0, 0), // $2 -> lastTime
		10,              // $3 -> limit
		"%%"}            // $4 -> title
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
	if param.Title != "" {
		args[3] = "%" + param.Title + "%"
	}

	fmt.Println(args)
	var rows []struct {
		Id          int       `db:"id"`
		Name        string    `db:"name"`
		Description string    `db:"description"`
		CreatedAt   time.Time `db:"created_at"`
	}
	if err := p.db.Select(&rows, query, args...); err != nil {
		return []entity.SerieListPage{}, fmt.Errorf(
			"persistence<Pg.Series>: %w", err)
	}

	var serieList []entity.SerieListPage
	var last *entity.SerieListPage
	for _, r := range rows {
		if last == nil || last.Id != r.Id {
			sl := entity.SerieListPage{
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
