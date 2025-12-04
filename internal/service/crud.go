package service

import (
	"fmt"
	"io"
	"os"

	"github.com/solsteace/misite/internal/entity"
)

func (s Service) InsertArticles(articles []entity.WriteArticle) error {
	contents := make([]string, len(articles))
	for idx, a := range articles {
		f, err := os.Open(a.Content)
		if err != nil {
			return fmt.Errorf("service<Service.InsertArticles>: %w", err)
		}

		content, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("service<Service.InsertArticles>: %w", err)
		}
		contents[idx] = string(content)
	}

	if err := s.store.InsertArticles(articles, contents); err != nil {
		return fmt.Errorf("service<Service.InsertArticles>: %w", err)
	}
	return nil
}

func (s Service) UpsertArticles(articles []entity.WriteArticle) error {
	contents := make([]string, len(articles))
	for idx, a := range articles {
		f, err := os.Open(a.Content)
		if err != nil {
			return fmt.Errorf("service<Service.UpsertArticles>: %w", err)
		}

		content, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("service<Service.UpsertArticles>: %w", err)
		}
		contents[idx] = string(content)
	}

	if err := s.store.UpsertArticles(articles, contents); err != nil {
		return fmt.Errorf("service<Service.UpsertArticles>: %w", err)
	}
	return nil
}

func (s Service) DeleteArticles(articles []entity.DeleteById) error {
	if err := s.store.DeleteArticles(articles); err != nil {
		return fmt.Errorf("service<Service.DeleteArticles>: %w", err)
	}
	return nil
}

func (s Service) InsertArticleTags(articleTags []entity.WriteArticleTag) error {
	if err := s.store.InsertArticlesTags(articleTags); err != nil {
		return fmt.Errorf("service<Service.InsertArticleTags>: %w", err)
	}
	return nil
}

func (s Service) UpsertArticleTags(articleTags []entity.WriteArticleTag) error {
	if err := s.store.UpsertArticleTags(articleTags); err != nil {
		return fmt.Errorf("service<Service.UpsertArticleTags>: %w", err)
	}
	return nil
}

func (s Service) DeleteArticleTags(articleTags []entity.DeleteById) error {
	if err := s.store.DeleteArticleTags(articleTags); err != nil {
		return fmt.Errorf("service<Service.DeleteArticleTags>: %w", err)
	}
	return nil
}

func (s Service) InsertProjects(projects []entity.WriteProject) error {
	contents := make([]string, len(projects))
	for idx, a := range projects {
		f, err := os.Open(a.Description)
		if err != nil {
			return fmt.Errorf("service<Service.InsertProjects>: %w", err)
		}

		content, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("service<Service.InsertProjects>: %w", err)
		}
		contents[idx] = string(content)
	}

	if err := s.store.InsertProjects(projects, contents); err != nil {
		return fmt.Errorf("service<Service.InsertProjects>: %w", err)
	}
	return nil
}

func (s Service) UpsertProjects(projects []entity.WriteProject) error {
	contents := make([]string, len(projects))
	for idx, a := range projects {
		f, err := os.Open(a.Description)
		if err != nil {
			return fmt.Errorf("service<Service.UpsertArticles>: %w", err)
		}

		content, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("service<Service.UpsertArticles>: %w", err)
		}
		contents[idx] = string(content)
	}

	if err := s.store.UpsertProjects(projects, contents); err != nil {
		return fmt.Errorf("service<Service.UpsertProjects>: %w", err)
	}
	return nil
}

func (s Service) DeleteProjects(projects []entity.DeleteById) error {
	if err := s.store.DeleteProjects(projects); err != nil {
		return fmt.Errorf("service<Service.DeleteProjects>: %w", err)
	}
	return nil
}

func (s Service) InsertProjectTags(projectTags []entity.WriteProjectTag) error {
	if err := s.store.InsertProjectTags(projectTags); err != nil {
		return fmt.Errorf("service<Service.InsertProjectTags>: %w", err)
	}
	return nil
}

func (s Service) UpsertProjectTags(projectTags []entity.WriteProjectTag) error {
	if err := s.store.UpsertProjectTags(projectTags); err != nil {
		return fmt.Errorf("service<Service.UpsertProjectTags>: %w", err)
	}
	return nil
}

func (s Service) DeleteProjectTags(projectTags []entity.DeleteById) error {
	if err := s.store.DeleteProjectTags(projectTags); err != nil {
		return fmt.Errorf("service<Service.DeleteProjectTags>: %w", err)
	}
	return nil
}

func (s Service) InsertProjectLinks(projectLinks []entity.WriteProjectLink) error {
	if err := s.store.InsertProjectLinks(projectLinks); err != nil {
		return fmt.Errorf("service<Service.InsertProjectLinks>: %w", err)
	}
	return nil
}

func (s Service) UpsertProjectLinks(projectLinks []entity.WriteProjectLink) error {
	if err := s.store.UpsertProjectLinks(projectLinks); err != nil {
		return fmt.Errorf("service<Service.UpsertProjectLinks>: %w", err)
	}
	return nil
}

func (s Service) DeleteProjectLinks(projectLinks []entity.DeleteById) error {
	if err := s.store.DeleteProjectLinks(projectLinks); err != nil {
		return fmt.Errorf("service<Service.DeleteProjectLinks>: %w", err)
	}
	return nil
}

func (s Service) InsertTags(tags []entity.WriteTag) error {
	if err := s.store.InsertTags(tags); err != nil {
		return fmt.Errorf("service<Service.UpsertTags>: %w", err)
	}
	return nil
}

func (s Service) UpsertTags(tags []entity.WriteTag) error {
	if err := s.store.UpsertTags(tags); err != nil {
		return fmt.Errorf("service<Service.UpsertTags>: %w", err)
	}
	return nil
}

func (s Service) DeleteTags(tags []entity.DeleteById) error {
	if err := s.store.DeleteTags(tags); err != nil {
		return fmt.Errorf("service<Service.DeleteTags>: %w", err)
	}
	return nil
}

func (s Service) InsertSeries(series []entity.WriteSerie) error {
	if err := s.store.InsertSeries(series); err != nil {
		return fmt.Errorf("service<Service.InsertSeries>: %w", err)
	}
	return nil
}

func (s Service) UpsertSeries(series []entity.WriteSerie) error {
	if err := s.store.UpsertSeries(series); err != nil {
		return fmt.Errorf("service<Service.UpsertSeries>: %w", err)
	}
	return nil
}

func (s Service) DeleteSeries(series []entity.DeleteById) error {
	if err := s.store.DeleteSeries(series); err != nil {
		return fmt.Errorf("service<Service.DeleteSeries>: %w", err)
	}
	return nil
}
