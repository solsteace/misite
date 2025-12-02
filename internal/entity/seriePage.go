package entity

import "time"

// The model for viewing `serie` entry in `serie` page
type SeriePage struct {
	Id          int
	NArticle    int
	NProject    int
	Name        string
	Thumbnail   string
	Description string
}

// Articles associated with the serie sorted in
// ascending order by their appearance on the serie
type SeriePageArticleList struct {
	Id        int
	Title     string
	Synopsis  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Project associated with the serie
type SeriePageProjectList struct {
	Id        int
	Name      string
	Synopsis  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// The model for viewing `serie` entry in `serie_list` page
type SerieListPage struct {
	Id          int
	Name        string
	Description string
	CreatedAt   time.Time
}

// A serie entry is considered new for 5 days after its initial creation
func (sl SerieListPage) IsNew() bool {
	return time.Since(sl.CreatedAt) < time.Hour*24*5
}
