package persistence

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/solsteace/misite/internal/entity"
)

func (p Pg) InsertArticles(articles []entity.WriteArticle, contents []string) error {
	query := `
		INSERT INTO articles(
			title,
			subtitle,
			content)
		VALUES(
			:title,
			:subtitle,
			:content)`
	rows := make([]any, len(articles))
	for idx, a := range articles {
		rows[idx] = struct {
			Title    string `db:"title"`
			Subtitle string `db:"subtitle"`
			Content  string `db:"content"`
		}{
			Title:    a.Title,
			Subtitle: a.Subtitle,
			Content:  contents[idx]}
	}
	if _, err := p.db.NamedExec(query, rows); err != nil {
		return fmt.Errorf("persistence<Pg.InsertArticles>: %w", err)
	}
	return nil
}

func (p Pg) UpsertArticles(articles []entity.WriteArticle, contents []string) error {
	query := `
		INSERT INTO articles(
			id,
			title,
			subtitle,
			content)
		VALUES(
			:id,
			:title,
			:subtitle,
			:content)
		ON CONFLICT(id)
		DO UPDATE SET
			title = EXCLUDED.title,
			subtitle = EXCLUDED.subtitle,
			content = EXCLUDED.content`
	rows := make([]any, len(articles))
	for idx, a := range articles {
		rows[idx] = struct {
			Id       int    `db:"id"`
			Title    string `db:"title"`
			Subtitle string `db:"subtitle"`
			Content  string `db:"content"`
		}{
			Id:       a.Id,
			Title:    a.Title,
			Subtitle: a.Subtitle,
			Content:  contents[idx]}
	}
	if _, err := p.db.NamedExec(query, rows); err != nil {
		return fmt.Errorf("persistence<Pg.UpsertArticles>: %w", err)
	}
	return nil
}

func (p Pg) DeleteArticles(articles []entity.DeleteById) error {
	targets := make([]any, len(articles))
	for idx, a := range articles {
		targets[idx] = a.Id
	}

	query, args, err := sqlx.In(
		`DELETE FROM articles WHERE id IN (?)`,
		targets)
	if err != nil {
		return fmt.Errorf("persistence<Pg.DeleteArticles>: %w", err)
	}

	if _, err := p.db.Exec(p.db.Rebind(query), args...); err != nil {
		return fmt.Errorf("persistence<Pg.DeleteArticles>: %w", err)
	}
	return nil
}

func (p Pg) InsertArticlesTags(articleTags []entity.WriteArticleTag) error {
	query := `
		INSERT INTO article_tags(
			article_id,
			tag_id)
		VALUES (
			:article_id,
			:tag_id)`
	rows := make([]any, len(articleTags))
	for idx, at := range articleTags {
		rows[idx] = struct {
			ArticleId int `db:"article_id"`
			TagId     int `db:"tag_id"`
		}{
			ArticleId: at.ArticleId,
			TagId:     at.TagId}
	}
	if _, err := p.db.NamedExec(query, rows); err != nil {
		return fmt.Errorf("persistence<Pg.InsertArticleTags>: %w", err)
	}
	return nil
}

func (p Pg) UpsertArticleTags(articleTags []entity.WriteArticleTag) error {
	query := `
		INSERT INTO article_tags(
			id,
			article_id,
			tag_id)
		VALUES (
			:id,
			:article_id,
			:tag_id)
		ON CONFLICT(id)
		DO UPDATE SET
			article_id = EXCLUDED.article_id,
			tag_id = EXCLUDED.tag_id`
	rows := make([]any, len(articleTags))
	for idx, at := range articleTags {
		rows[idx] = struct {
			Id        int `db:"id"`
			ArticleId int `db:"article_id"`
			TagId     int `db:"tag_id"`
		}{
			Id:        at.Id,
			ArticleId: at.ArticleId,
			TagId:     at.TagId}
	}
	if _, err := p.db.NamedExec(query, rows); err != nil {
		return fmt.Errorf("persistence<Pg.UpsertArticleTags>: %w", err)
	}
	return nil
}

func (p Pg) DeleteArticleTags(articleTags []entity.DeleteById) error {
	targets := make([]any, len(articleTags))
	for idx, at := range articleTags {
		targets[idx] = at.Id
	}

	query, args, err := sqlx.In(
		`DELETE FROM articles_tags WHERE id IN (?)`,
		targets...)
	if err != nil {
		return fmt.Errorf("persistence<Pg.DeleteArticleTags>: %w", err)
	}

	if _, err := p.db.Exec(p.db.Rebind(query), args...); err != nil {
		return fmt.Errorf("persistence<Pg.DeleteArticleTags>: %w", err)
	}
	return nil
}

func (p Pg) InsertProjects(projects []entity.WriteProject, contents []string) error {
	query := `
		INSERT INTO projects(
			name, 
			synopsis,
			description)
		VALUES(
			:name,
			:synopsis,
			:description)`
	rows := make([]any, len(projects))
	for idx, p := range projects {
		rows[idx] = struct {
			Name        string `db:"name"`
			Synopsis    string `db:"synopsis"`
			Description string `db:"description"`
		}{
			Name:        p.Name,
			Synopsis:    p.Synopsis,
			Description: contents[idx]}
	}
	if _, err := p.db.NamedExec(query, rows); err != nil {
		return fmt.Errorf("persistence<Pg.InsertProjects>: %w", err)
	}
	return nil
}

func (p Pg) UpsertProjects(projects []entity.WriteProject, contents []string) error {
	query := `
		INSERT INTO projects(
			id,
			name, 
			synopsis,
			description)
		VALUES(
			:id,
			:name,
			:synopsis,
			:description)
		ON CONFLICT(id)
		DO UPDATE SET
			name = EXCLUDED.name,
			synopsis = EXCLUDED.synopsis,
			description = EXCLUDED.description`
	rows := make([]any, len(projects))
	for idx, p := range projects {
		rows[idx] = struct {
			Id          int    `db:"id"`
			Name        string `db:"name"`
			Synopsis    string `db:"synopsis"`
			Description string `db:"description"`
		}{
			Id:          p.Id,
			Name:        p.Name,
			Synopsis:    p.Synopsis,
			Description: contents[idx]}
	}
	if _, err := p.db.NamedExec(query, rows); err != nil {
		return fmt.Errorf("persistence<Pg.UpsertProjects>: %w", err)
	}
	return nil
}

func (p Pg) DeleteProjects(projects []entity.DeleteById) error {
	targets := make([]any, len(projects))
	for idx, a := range projects {
		targets[idx] = a.Id
	}

	query, args, err := sqlx.In(
		`DELETE FROM projects WHERE id IN (?)`,
		targets)
	if err != nil {
		return fmt.Errorf("persistence<Pg.DeleteProjects>: %w", err)
	}

	if _, err := p.db.Exec(p.db.Rebind(query), args...); err != nil {
		return fmt.Errorf("persistence<Pg.DeleteProjects>: %w", err)
	}
	return nil
}

func (p Pg) InsertProjectTags(projectTags []entity.WriteProjectTag) error {
	query := `
		INSERT INTO project_tags(
			project_id, 
			tag_id)
		VALUES(
			:project_id,
			:tag_id)`
	rows := make([]any, len(projectTags))
	for idx, pt := range projectTags {
		rows[idx] = struct {
			ProjectId int `db:"project_id"`
			TagId     int `db:"tag_id"`
		}{
			ProjectId: pt.ProjectId,
			TagId:     pt.TagId}
	}
	if _, err := p.db.NamedExec(query, rows); err != nil {
		return fmt.Errorf("persistence<Pg.InsertProjectTags>: %w", err)
	}
	return nil
}

func (p Pg) UpsertProjectTags(projectTags []entity.WriteProjectTag) error {
	query := `
		INSERT INTO project_tags(
			id,
			project_id, 
			tag_id)
		VALUES(
			:id,
			:project_id,
			:tag_id)
		ON CONFLICT(id)
		DO UPDATE SET
			project_id = EXCLUDED.project_id,
			tag_id = EXCLUDED.tag_id`
	rows := make([]any, len(projectTags))
	for idx, pt := range projectTags {
		rows[idx] = struct {
			Id        int `db:"id"`
			ProjectId int `db:"project_id"`
			TagId     int `db:"tag_id"`
		}{
			Id:        pt.Id,
			ProjectId: pt.ProjectId,
			TagId:     pt.TagId}
	}
	if _, err := p.db.NamedExec(query, rows); err != nil {
		return fmt.Errorf("persistence<Pg.UpsertProjectTags>: %w", err)
	}
	return nil
}

func (p Pg) DeleteProjectTags(projectTags []entity.DeleteById) error {
	targets := make([]any, len(projectTags))
	for idx, a := range projectTags {
		targets[idx] = a.Id
	}

	query, args, err := sqlx.In(
		`DELETE FROM project_tags WHERE id IN (?)`,
		targets)
	if err != nil {
		return fmt.Errorf("persistence<Pg.DeleteProjects>: %w", err)
	}

	if _, err := p.db.Exec(p.db.Rebind(query), args...); err != nil {
		return fmt.Errorf("persistence<Pg.DeleteProjects>: %w", err)
	}
	return nil
}

func (p Pg) InsertProjectLinks(projectLinks []entity.WriteProjectLink) error {
	query := `
		INSERT INTO project_links(
			project_id, 
			display_text,
			url)
		VALUES(
			:project_id,
			:display_text,
			:url)`
	rows := make([]any, len(projectLinks))
	for idx, pl := range projectLinks {
		rows[idx] = struct {
			ProjectId   int    `db:"project_id"`
			DisplayText string `db:"display_text"`
			Url         string `db:"url"`
		}{
			ProjectId:   pl.ProjectId,
			DisplayText: pl.DisplayText,
			Url:         pl.Url}
	}
	if _, err := p.db.NamedExec(query, rows); err != nil {
		return fmt.Errorf("persistence<Pg.InsertProjectLinks>: %w", err)
	}
	return nil
}

func (p Pg) UpsertProjectLinks(projectLinks []entity.WriteProjectLink) error {
	query := `
		INSERT INTO project_links(
			id,
			project_id, 
			display_text,
			url)
		VALUES(
			:id,
			:project_id,
			:display_text,
			:url)
		ON CONFLICT(id)
		DO UPDATE SET
			project_id = EXCLUDED.project_id,
			display_text = EXCLUDED.display_text,
			url = EXCLUDED.url`
	rows := make([]any, len(projectLinks))
	for idx, pl := range projectLinks {
		rows[idx] = struct {
			Id          int    `db:"id"`
			ProjectId   int    `db:"project_id"`
			DisplayText string `db:"display_text"`
			Url         string `db:"url"`
		}{
			Id:          pl.Id,
			ProjectId:   pl.ProjectId,
			DisplayText: pl.DisplayText,
			Url:         pl.Url}
	}
	if _, err := p.db.NamedExec(query, rows); err != nil {
		return fmt.Errorf("persistence<Pg.UpsertProjectLinks>: %w", err)
	}
	return nil
}

func (p Pg) DeleteProjectLinks(projectLinks []entity.DeleteById) error {
	targets := make([]any, len(projectLinks))
	for idx, a := range projectLinks {
		targets[idx] = a.Id
	}

	query, args, err := sqlx.In(
		`DELETE FROM project_links WHERE id IN (?)`,
		targets)
	if err != nil {
		return fmt.Errorf("persistence<Pg.DeleteProjectLinks>: %w", err)
	}

	if _, err := p.db.Exec(p.db.Rebind(query), args...); err != nil {
		return fmt.Errorf("persistence<Pg.DeleteProjectLinks>: %w", err)
	}
	return nil
}

func (p Pg) InsertTags(tags []entity.WriteTag) error {
	query := `
		INSERT INTO tags(name)
		VALUES(:name)`
	rows := make([]any, len(tags))
	for idx, t := range tags {
		rows[idx] = struct {
			Name string `db:"name"`
		}{Name: t.Name}
	}
	if _, err := p.db.NamedExec(query, rows); err != nil {
		return fmt.Errorf("persistence<Pg.InsertTags>: %w", err)
	}
	return nil
}

func (p Pg) UpsertTags(tags []entity.WriteTag) error {
	query := `
		INSERT INTO tags(id, name)
		VALUES(:id, :name)
		ON CONFLICT(id)
		DO UPDATE SET
			name = EXCLUDED.name `
	rows := make([]any, len(tags))
	for idx, t := range tags {
		rows[idx] = struct {
			Id   int    `db:"id"`
			Name string `db:"name"`
		}{
			Id:   t.Id,
			Name: t.Name}
	}
	if _, err := p.db.NamedExec(query, rows); err != nil {
		return fmt.Errorf("persistence<Pg.UpsertTags>: %w", err)
	}
	return nil
}

func (p Pg) DeleteTags(tags []entity.DeleteById) error {
	targets := make([]any, len(tags))
	for idx, a := range tags {
		targets[idx] = a.Id
	}

	query, args, err := sqlx.In(
		`DELETE FROM tags WHERE id IN (?)`,
		targets)
	if err != nil {
		return fmt.Errorf("persistence<Pg.DeleteTags>: %w", err)
	}

	if _, err := p.db.Exec(p.db.Rebind(query), args...); err != nil {
		return fmt.Errorf("persistence<Pg.DeleteTags>: %w", err)
	}
	return nil
}

func (p Pg) InsertSeries(series []entity.WriteSerie) error {
	query := `
		INSERT INTO series(
			name,
			thumbnail,
			description)
		VALUES(
			:name,
			:thumbnail,
			:description)`
	rows := make([]any, len(series))
	for idx, s := range series {
		rows[idx] = struct {
			Name        string `db:"name"`
			Thumbnail   string `db:"thumbnail"`
			Description string `db:"description"`
		}{
			Name:        s.Name,
			Thumbnail:   s.Thumbnail,
			Description: s.Description}
	}
	if _, err := p.db.NamedExec(query, rows); err != nil {
		return fmt.Errorf("persistence<Pg.InsertSeries>: %w", err)
	}
	return nil
}

func (p Pg) UpsertSeries(series []entity.WriteSerie) error {
	query := `
		INSERT INTO series(
			id,
			name,
			thumbnail,
			description)
		VALUES(
			:id,
			:name,
			:thumbnail,
			:description)
		ON CONFLICT(id)
		DO UPDATE SET
			name = EXCLUDED.name,
			thumbnail = EXCLUDED.thumbnail,
			description = EXCLUDED.decsription `
	rows := make([]any, len(series))
	for idx, s := range series {
		rows[idx] = struct {
			Id          int    `db:"id"`
			Name        string `db:"name"`
			Thumbnail   string `db:"thumbnail"`
			Description string `db:"description"`
		}{
			Id:          s.Id,
			Name:        s.Name,
			Thumbnail:   s.Thumbnail,
			Description: s.Description}
	}
	if _, err := p.db.NamedExec(query, rows); err != nil {
		return fmt.Errorf("persistence<Pg.UpsertSeries>: %w", err)
	}
	return nil
}

func (p Pg) DeleteSeries(series []entity.DeleteById) error {
	targets := make([]any, len(series))
	for idx, a := range series {
		targets[idx] = a.Id
	}

	query, args, err := sqlx.In(
		`DELETE FROM series WHERE id IN (?)`,
		targets)
	if err != nil {
		return fmt.Errorf("persistence<Pg.DeleteSeries>: %w", err)
	}

	if _, err := p.db.Exec(p.db.Rebind(query), args...); err != nil {
		return fmt.Errorf("persistence<Pg.DeleteSeries>: %w", err)
	}
	return nil
}
