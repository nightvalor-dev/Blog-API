package comment

type CreateCommentRequest struct {
	Content string `json:"content"`
	BlogId  int    `json:"blog_id"`
	UserId  int    `json:"user_id"`
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}

type CommentResponse struct {
	Content string `json:"content"`
	BlogId  int    `json:"blog_id"`
	UserId  int    `json:"user_id"`
}
