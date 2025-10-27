package entity

import "time"

// The model for viewing `serie` entry in `serie` page
type Serie2 struct {
	Id          int
	Name        string
	Thumbnail   string
	Description string

	// Articles associated with the serie sorted in
	// ascending order by their appearance on the serie
	Article []struct {
		Id        int
		Title     string
		Synopsis  string
		Thumbnail string
		Order     int
		CreatedAt time.Time
	}

	// Project associated with the serie
	Project []struct {
		Id        int
		Name      string
		Thumbnail string
		Synopsis  string
	}
}

// The model for viewing `serie` entry in `serie_list` page
type SerieList struct {
	Id          int
	Name        string
	Thumbnail   string
	Description string
}
