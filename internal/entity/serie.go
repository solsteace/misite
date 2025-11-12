package entity

import "time"

// The model for viewing `serie` entry in `serie` page
type Serie struct {
	Id          int
	NArticle    int
	NProject    int
	Name        string
	Thumbnail   string
	Description string
}

// Articles associated with the serie sorted in
// ascending order by their appearance on the serie
type SerieArticleList struct {
	Id        int
	Title     string
	Synopsis  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Project associated with the serie
type SerieProjectList struct {
	Id        int
	Name      string
	Synopsis  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// The model for viewing `serie` entry in `serie_list` page
type SerieList struct {
	Id          int
	Name        string
	Description string
	CreatedAt   time.Time
}

// A serie entry is considered new for 5 days after its initial creation
func (sl SerieList) IsNew() bool {
	return time.Now().Sub(sl.CreatedAt) < time.Hour*24*5
}
