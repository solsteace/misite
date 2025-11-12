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
	CreatedAt time.Time
	UpdatedAt time.Time

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

func (a Article) DisplayTime() string {
	var timestamp string
	if diff := time.Now().Sub(a.CreatedAt); diff < time.Hour*24 {
		timestamp += "Today"
	} else if diff < time.Hour*24*14 {
		timestamp = fmt.Sprintf("%.0f days ago", diff.Hours()/24)
	} else if diff < time.Hour*24*365 {
		timestamp = a.CreatedAt.Format("Jan 02")
	} else {
		timestamp = a.CreatedAt.Format("Jan 02, 2006")
	}

	if a.UpdatedAt == a.CreatedAt {
		return timestamp
	}

	timestamp += " "
	if diff := time.Now().Sub(a.UpdatedAt); diff < time.Hour*24 {
		timestamp += "(updated today)"
	} else if diff < time.Hour*24*14 {
		timestamp += fmt.Sprintf("(updated %.0f days ago)", diff.Hours()/24)
	} else if diff < time.Hour*24*365 {
		timestamp += fmt.Sprintf("(updated at %s)", a.UpdatedAt.Format("Jan 02"))
	} else {
		timestamp += fmt.Sprintf("(updated at %s)", a.UpdatedAt.Format("Jan 02, 2006"))
	}
	return timestamp
}

type ArticleList struct {
	Id        int
	Title     string
	Subtitle  string
	CreatedAt time.Time
	UpdatedAt time.Time

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
func ArticleIsNew(createdAt time.Time) bool {
	return time.Now().Sub(createdAt) < time.Hour*24*5
}

// An article entry is considered recently updated for 3 days after its latest change
func ArticleIsRecentlyUpdated(createdAt, updatedAt time.Time) bool {
	return updatedAt != createdAt && time.Now().Sub(updatedAt) < time.Hour*24*3
}
