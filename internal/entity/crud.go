package entity

type WriteArticle struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"` // path to HTML file containing the content
}

type WriteArticleTag struct {
	Id        int `json:"id"`
	ArticleId int `json:"article_id"`
	TagId     int `json:"tag_id"`
}

type WriteProject struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Synopsis    string `json:"synopsis"`
	Description string `json:"description"` // path to HTML file containing the content
}

type WriteProjectTag struct {
	Id        int `json:"id"`
	ProjectId int `json:"project_id"`
	TagId     int `json:"tag_id"`
}

type WriteProjectLink struct {
	Id          int    `json:"id"`
	ProjectId   int    `json:"project_id"`
	DisplayText string `json:"display_text"`
	Url         string `json:"url"`
}

type WriteTag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type WriteSerie struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Thumbnail   string `json:"thumbnail"`
	Description string `json:"description"`
}

type DeleteById struct {
	Id int `json:"id"`
}
