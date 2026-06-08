package media

type MediaRequest struct {
	BlogId int `json:"blog_id"`
}

type MediaResponse struct {
	MediaId int    `json:"media_id"`
	BlogId  int    `json:"blog_id"`
	Url     string `json:"url"`
}
