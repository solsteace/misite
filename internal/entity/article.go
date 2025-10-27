package entity

import "time"

// The model to show an article on its specification page
type Article struct {
	Id        int
	Title     string
	Subtitle  string
	Content   string
	Thumbnail string
	CreatedAt time.Time

	// an article series that accompanies the project, if any (some kind of devblog, if you will)
	Serie *struct {
		Id   int
		Name string
	}

	// the associated tags
	Tag []struct {
		Id   int
		Name string
	}
}

type ArticleList struct {
	Id        int
	Title     string
	Subtitle  string
	Thumbnail string
	CreatedAt time.Time

	// an article series that accompanies the project, if any (some kind of devblog, if you will)
	Serie *struct {
		Id   int
		Name string
	}

	// the associated tags
	Tag []struct {
		Id   int
		Name string
	}
}
