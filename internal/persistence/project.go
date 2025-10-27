package persistence

import (
	"database/sql"
	"fmt"

	"github.com/solsteace/misite/internal/entity"
	"github.com/solsteace/misite/internal/utility/lib/oops"
)

type ProjectsQueryParam struct {
	Page  int
	Limit int
	TagId []int
}

func (p Pg) Projects(param ProjectsQueryParam) ([]entity.ProjectList, error) {
	query := `
		WITH matching_projects_by_tag AS(
			SELECT 
				project_id, 
				COUNT(DISTINCT tag_id) AS "n_tag"
			FROM project_tags
			WHERE 
				($1::int[] IS NULL)
				OR tag_id = ANY($1)
			GROUP BY project_id
			HAVING
				($1::int[] IS NULL)
				OR COUNT(tag_id) = CARDINALITY($1)
			LIMIT $2 OFFSET $3)
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
		WHERE
			(SELECT true
			FROM matching_projects_by_tag
			WHERE matching_projects_by_tag.project_id = projects.id)
		ORDER BY projects.id`
	if param.Limit < 1 {
		param.Limit = 10
	}
	if param.Page < 1 {
		param.Page = 1
	}
	args := []any{
		nil,
		param.Limit,
		(param.Page - 1) * param.Limit}
	if len(param.TagId) > 0 {
		args[0] = param.TagId
	}

	var rows []struct {
		Id        int    `db:"id"`
		Name      string `db:"name"`
		Thumbnail string `db:"thumbnail"`
		Synopsis  string `db:"synopsis"`

		Tag struct {
			Id   sql.Null[int]    `db:"id"`
			Name sql.Null[string] `db:"name"`
		}
	}
	if err := p.db.Select(&rows, query, args...); err != nil {
		return []entity.ProjectList{}, fmt.Errorf(
			"persistence<Pg.Projects>: %w", err)
	} else if len(rows) == 0 {
		return []entity.ProjectList{}, nil
	}

	var projects []entity.ProjectList
	var lastProject *entity.ProjectList
	var insertedTag map[int]struct{}
	for _, r := range rows {
		if lastProject == nil || lastProject.Id != r.Id {
			insertedTag = map[int]struct{}{}
			projects = append(projects, entity.ProjectList{
				Id:        r.Id,
				Name:      r.Name,
				Thumbnail: r.Thumbnail,
				Synopsis:  r.Synopsis})
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
	}
	return projects, nil
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

	var rows []struct {
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
		ORDER BY tag_count.count DESC`
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
