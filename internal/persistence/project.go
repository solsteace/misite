package persistence

import (
	"database/sql"
	"fmt"

	"github.com/solsteace/misite/internal/entity"
	"github.com/solsteace/misite/internal/utility/lib/oops"
)

type ProjectsQueryParam struct {
	Page    int
	Limit   int
	TagId   []int
	SerieId []int
}

type pgProjects struct {
	Id        int    `db:"id"`
	Name      string `db:"name"`
	Thumbnail string `db:"thumbnail"`
	Synopsis  string `db:"synopsis"`

	Tag struct {
		Id   sql.Null[int]    `db:"id"`
		Name sql.Null[string] `db:"name"`
	}
}

func (p Pg) Projects(param ProjectsQueryParam) ([]entity.Project, error) {
	if param.Limit < 1 {
		param.Limit = 10
	}
	if param.Page < 1 {
		param.Page = 1
	}
	query := `
		SELECT
			projects.id AS "id",
			projects.name AS "name",
			projects.thumbnail AS "thumbnail",
			projects.synopsis AS "synopsis",
			tags.id AS "tag.id",
			tags.name AS "tag.name"
		FROM projects
		LEFT JOIN project_tags ON project_tags.project_id = projects.id
		LEFT JOIN tags ON project_tags.tag_id = tags.id
		ORDER BY projects.id
		LIMIT $1
		OFFSET $2`
	args := []any{
		param.Limit,
		(param.Page - 1) * param.Limit}
	var rows []pgProjects
	if err := p.db.Select(&rows, query, args...); err != nil {
		return []entity.Project{}, fmt.Errorf(
			"persistence<Pg.Projects>: %w", err)
	} else if len(rows) == 0 {
		return []entity.Project{}, nil
	}

	var insertedTag map[int]struct{}
	var projects []entity.Project
	var lastProject *entity.Project
	for _, r := range rows {
		if lastProject == nil || lastProject.Id != r.Id {
			insertedTag = map[int]struct{}{}
			projects = append(projects, entity.Project{
				Id:        r.Id,
				Name:      r.Name,
				Thumbnail: r.Thumbnail,
				Synopsis:  r.Synopsis})
			lastProject = &projects[len(projects)-1]
		}
		if r.Tag.Id.Valid {
			if _, ok := insertedTag[r.Tag.Id.V]; !ok {
				insertedTag[r.Tag.Id.V] = struct{}{}
				lastProject.Tag = append(lastProject.Tag, entity.Tag{
					Id:   r.Tag.Id.V,
					Name: r.Tag.Name.V})
			}
		}
	}
	return projects, nil
}

type pgProject struct {
	Id          int    `db:"id"`
	Name        string `db:"name"`
	Thumbnail   string `db:"thumbnail"`
	Synopsis    string `db:"synopsis"`
	Description string `db:"description"`

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

func (p Pg) Project(id int) (entity.Project, error) {
	query := `
		SELECT
			projects.id AS "id",
			projects.name AS "name",
			projects.thumbnail AS "thumbnail",
			projects.synopsis AS "synopsis",
			projects.description AS "description",
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
	var rows []pgProject
	if err := p.db.Select(&rows, query, args...); err != nil {
		return entity.Project{}, fmt.Errorf(
			"persistence<pg.Project>: %w", err)
	} else if len(rows) == 0 {
		return entity.Project{}, fmt.Errorf(
			"persistence<pg.Project>: %w", oops.NotFound{})
	}

	projectRow := rows[0]
	project := entity.Project{
		Id:          projectRow.Id,
		Name:        projectRow.Name,
		Thumbnail:   projectRow.Thumbnail,
		Synopsis:    projectRow.Synopsis,
		Description: projectRow.Description}
	if projectRow.Serie.Id.Valid {
		project.Serie = &entity.Serie{
			Id:   projectRow.Serie.Id.V,
			Name: projectRow.Serie.Name.V}
	}

	insertedLink := map[int]struct{}{}
	insertedTag := map[int]struct{}{}
	for _, r := range rows {
		if r.Link.Id.Valid {
			if _, ok := insertedLink[r.Link.Id.V]; !ok {
				insertedLink[r.Link.Id.V] = struct{}{}
				project.Link = append(project.Link, entity.ProjectLink{
					Id:          r.Link.Id.V,
					DisplayText: r.Link.DisplayText.V,
					Url:         r.Link.Url.V})
			}
		}
		if r.Tag.Id.Valid {
			if _, ok := insertedTag[r.Tag.Id.V]; !ok {
				insertedTag[r.Tag.Id.V] = struct{}{}
				project.Tag = append(project.Tag, entity.Tag{
					Id:   r.Tag.Id.V,
					Name: r.Tag.Name.V})
			}
		}
	}
	return project, nil
}
