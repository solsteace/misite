package entity

// The model to show a project on its specification page
type Project struct {
	Id          int
	Name        string
	Thumbnail   string
	Synopsis    string
	Description string

	Serie *Serie        // an article series that accompanies the project, if any (some kind of devblog, if you will)
	Tag   []Tag         // the associated tags
	Link  []ProjectLink // the related links (deployment, references, etc.)
}
