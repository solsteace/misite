package controller

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/solsteace/misite/internal/entity"
)

func (c Controller) InsertArticles(f *os.File) error {
	var data struct {
		Articles []entity.WriteArticle `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.InsertArticle>: %w", err)
	}

	if err := c.service.InsertArticles(data.Articles); err != nil {
		return fmt.Errorf("controller<Controller.InsertArticle>: %w", err)
	}
	return nil
}

func (c Controller) UpsertArticles(f *os.File) error {
	var data struct {
		Articles []entity.WriteArticle `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.UpsertArticle>: %w", err)
	}

	if err := c.service.UpsertArticles(data.Articles); err != nil {
		return fmt.Errorf("controller<Controller.UpsertArticle>: %w", err)
	}
	return nil
}

func (c Controller) DeleteArticles(f *os.File) error {
	var data struct {
		Articles []entity.DeleteById `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.DeleteArticle>: %w", err)
	}

	if err := c.service.DeleteArticles(data.Articles); err != nil {
		return fmt.Errorf("controller<Controller.DeleteArticle>: %w", err)
	}
	return nil
}

func (c Controller) InsertArticleTags(f *os.File) error {
	var data struct {
		ArticleTags []entity.WriteArticleTag `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.InsertArticleTag>: %w", err)
	}

	if err := c.service.InsertArticleTags(data.ArticleTags); err != nil {
		return fmt.Errorf("controller<Controller.InsertArticleTag>: %w", err)
	}
	return nil
}

func (c Controller) UpsertArticleTags(f *os.File) error {
	var data struct {
		ArticleTags []entity.WriteArticleTag `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.UpsertArticleTag>: %w", err)
	}

	if err := c.service.UpsertArticleTags(data.ArticleTags); err != nil {
		return fmt.Errorf("controller<Controller.UpsertArticleTag>: %w", err)
	}
	return nil
}

func (c Controller) DeleteArticleTags(f *os.File) error {
	var data struct {
		ArticleTags []entity.DeleteById `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.DeleteArticleTag>: %w", err)
	}

	if err := c.service.DeleteArticleTags(data.ArticleTags); err != nil {
		return fmt.Errorf("controller<Controller.DeleteArticleTag>: %w", err)
	}
	return nil
}

func (c Controller) InsertProjects(f *os.File) error {
	var data struct {
		Projects []entity.WriteProject `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.InsertProject>: %w", err)
	}

	if err := c.service.InsertProjects(data.Projects); err != nil {
		return fmt.Errorf("controller<Controller.InsertProject>: %w", err)
	}
	return nil
}

func (c Controller) UpsertProjects(f *os.File) error {
	var data struct {
		Projects []entity.WriteProject `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.UpsertProject>: %w", err)
	}

	if err := c.service.UpsertProjects(data.Projects); err != nil {
		return fmt.Errorf("controller<Controller.UpsertProject>: %w", err)
	}
	return nil
}

func (c Controller) DeleteProjects(f *os.File) error {
	var data struct {
		Projects []entity.DeleteById `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.DeleteProject>: %w", err)
	}

	if err := c.service.DeleteProjects(data.Projects); err != nil {
		return fmt.Errorf("controller<Controller.DeleteProject>: %w", err)
	}
	return nil
}

func (c Controller) InsertProjectTags(f *os.File) error {
	var data struct {
		ProjectTags []entity.WriteProjectTag `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.InsertProjectTag>: %w", err)
	}

	if err := c.service.InsertProjectTags(data.ProjectTags); err != nil {
		return fmt.Errorf("controller<Controller.InsertProjectTag>: %w", err)
	}
	return nil
}

func (c Controller) UpsertProjectTags(f *os.File) error {
	var data struct {
		ProjectTags []entity.WriteProjectTag `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.UpsertProjectTag>: %w", err)
	}

	if err := c.service.UpsertProjectTags(data.ProjectTags); err != nil {
		return fmt.Errorf("controller<Controller.UpsertProjectTag>: %w", err)
	}
	return nil
}

func (c Controller) DeleteProjectTags(f *os.File) error {
	var data struct {
		ProjectTags []entity.DeleteById `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.DeleteProjectTag>: %w", err)
	}

	if err := c.service.DeleteProjectTags(data.ProjectTags); err != nil {
		return fmt.Errorf("controller<Controller.DeleteProjectTag>: %w", err)
	}
	return nil
}

func (c Controller) InsertProjectLinks(f *os.File) error {
	var data struct {
		ProjectLink []entity.WriteProjectLink `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.InsertProjectLink>: %w", err)
	}

	if err := c.service.InsertProjectLinks(data.ProjectLink); err != nil {
		return fmt.Errorf("controller<Controller.InsertProjectLink>: %w", err)
	}
	return nil
}

func (c Controller) UpsertProjectLinks(f *os.File) error {
	var data struct {
		ProjectLink []entity.WriteProjectLink `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.UpsertProjectLink>: %w", err)
	}

	if err := c.service.UpsertProjectLinks(data.ProjectLink); err != nil {
		return fmt.Errorf("controller<Controller.UpsertProjectLink>: %w", err)
	}
	return nil
}

func (c Controller) DeleteProjectLinks(f *os.File) error {
	var data struct {
		ProjectLink []entity.DeleteById `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.DeleteProjectLink>: %w", err)
	}

	if err := c.service.DeleteProjectTags(data.ProjectLink); err != nil {
		return fmt.Errorf("controller<Controller.DeleteProjectLink>: %w", err)
	}
	return nil
}

func (c Controller) InsertTags(f *os.File) error {
	var data struct {
		Tag []entity.WriteTag `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.InsertTag>: %w", err)
	}

	if err := c.service.InsertTags(data.Tag); err != nil {
		return fmt.Errorf("controller<Controller.InsertTag>: %w", err)
	}
	return nil
}

func (c Controller) UpsertTags(f *os.File) error {
	var data struct {
		Tag []entity.WriteTag `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.UpsertTag>: %w", err)
	}

	if err := c.service.UpsertTags(data.Tag); err != nil {
		return fmt.Errorf("controller<Controller.UpsertTag>: %w", err)
	}
	return nil
}

func (c Controller) DeleteTags(f *os.File) error {
	var data struct {
		Tag []entity.DeleteById `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.DeleteTag>: %w", err)
	}

	if err := c.service.DeleteProjectTags(data.Tag); err != nil {
		return fmt.Errorf("controller<Controller.DeleteTag>: %w", err)
	}
	return nil
}

func (c Controller) InsertSeries(f *os.File) error {
	var data struct {
		Serie []entity.WriteSerie `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.InsertSerie>: %w", err)
	}

	if err := c.service.InsertSeries(data.Serie); err != nil {
		return fmt.Errorf("controller<Controller.InsertSerie>: %w", err)
	}
	return nil
}

func (c Controller) UpsertSeries(f *os.File) error {
	var data struct {
		Serie []entity.WriteSerie `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.UpsertSerie>: %w", err)
	}

	if err := c.service.UpsertSeries(data.Serie); err != nil {
		return fmt.Errorf("controller<Controller.UpsertSerie>: %w", err)
	}
	return nil
}

func (c Controller) DeleteSeries(f *os.File) error {
	var data struct {
		Serie []entity.DeleteById `json:"data"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("controller<Controller.DeleteSerie>: %w", err)
	}

	if err := c.service.DeleteSeries(data.Serie); err != nil {
		return fmt.Errorf("controller<Controller.DeleteSerie>: %w", err)
	}
	return nil
}
