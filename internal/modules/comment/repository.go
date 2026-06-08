package comment

type CommentRepository interface {
	Create(comment CreateCommentRequest) error
	GetAll() ([]CommentResponse, error)
	GetById(id int) (CommentResponse, error)
	Update(id int, newComment UpdateCommentRequest) error
	Delete(id int) error
}
