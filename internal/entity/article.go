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

	Serie *Serie // an article series that accompanies the project, if any (some kind of devblog, if you will)
	Tag   []Tag  // the associated tags
}
