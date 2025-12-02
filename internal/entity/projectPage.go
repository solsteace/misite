package entity

import (
	"fmt"
	"time"
)

// The model to show a project on its specification page
type ProjectPage struct {
	Id          int
	Name        string
	Synopsis    string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// an article serie that accompanies the project, if any (some kind of devblog, if you will)
	Serie *struct {
		Id   int
		Name string
	}
	// the associated tags
	Tag []struct {
		Id   int
		Name string
	}
	// the related links (deployment, references, etc.)
	Link []struct {
		Id          int
		DisplayText string
		Url         string
	}
}

func (p ProjectPage) DisplayTime() string {
	var timestamp string
	if diff := time.Since(p.CreatedAt); diff < time.Hour*24 {
		timestamp += "Today"
	} else if diff < time.Hour*24*14 {
		timestamp = fmt.Sprintf("%.0f days ago", diff.Hours()/24)
	} else if diff < time.Hour*24*365 {
		timestamp = p.CreatedAt.Format("Jan 02")
	} else {
		timestamp = p.CreatedAt.Format("Jan 02, 2006")
	}

	if p.UpdatedAt.Equal(p.CreatedAt) {
		return timestamp
	}

	timestamp += " "
	if diff := time.Since(p.UpdatedAt); diff < time.Hour*24 {
		timestamp += "(updated today)"
	} else if diff < time.Hour*24*14 {
		timestamp += fmt.Sprintf("(updated %.0f days ago)", diff.Hours()/24)
	} else if diff < time.Hour*24*365 {
		timestamp += fmt.Sprintf("(updated at %s)", p.UpdatedAt.Format("Jan 02"))
	} else {
		timestamp += fmt.Sprintf("(updated at %s)", p.UpdatedAt.Format("Jan 02, 2006"))
	}
	return timestamp

}

type ProjectListPage struct {
	Id        int
	Name      string
	Synopsis  string
	CreatedAt time.Time
	UpdatedAt time.Time

	// an article serie that accompanies the project, if any (some kind of devblog, if you will)
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

// An project entry is considered new for 5 days after its initial creation
func ProjectIsNew(createdAt time.Time) bool {
	return time.Since(createdAt) < time.Hour*24*5
}

// An project entry is considered recently updated for 3 days after its latest change
func ProjectIsRecentlyUpdated(createdAt, updatedAt time.Time) bool {
	return !updatedAt.Equal(createdAt) && time.Since(updatedAt) < time.Hour*24*3
}
