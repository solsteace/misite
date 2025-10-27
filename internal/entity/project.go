package entity

// The model to show a project on its specification page
type Project struct {
	Id          int
	Name        string
	Thumbnail   string
	Synopsis    string
	Description string

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

type ProjectList struct {
	Id        int
	Name      string
	Thumbnail string
	Synopsis  string

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
