package entity

import (
	"fmt"
	"time"
)

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

func (a Article) DisplayCreationTime() string {
	if diff := time.Now().Sub(a.CreatedAt); diff < time.Hour*24*14 {
		return fmt.Sprintf("%.0f days ago", diff.Hours()/24)
	} else if diff < time.Hour*24*365 {
		return a.CreatedAt.Format("Jan 02")
	} else {
		return a.CreatedAt.Format("Jan 02, 2006")
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

// An article entry is considered new for 5 days after its initial creation
func (al ArticleList) IsNew() bool {
	return time.Now().Sub(al.CreatedAt) < time.Hour*24*365*2
}
