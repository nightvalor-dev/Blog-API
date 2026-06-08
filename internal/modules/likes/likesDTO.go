package likes

type LikeRequest struct {
	BlogId int `json:"blog_id"`
	UserId int `json:"user_id"`
}

type LikeResponse struct {
	BlogId int `json:"blog_id"`
	Count  int `json:"count"`
}
