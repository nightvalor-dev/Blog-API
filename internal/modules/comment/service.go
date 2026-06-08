package comment

type CommentService struct {
	repo CommentRepository
}

func NewCommentService(repo CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) Create(req CreateCommentRequest) error {
	return s.repo.Create(req)
}

func (s *CommentService) GetAll() ([]CommentResponse, error) {
	comments, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *CommentService) GetById(id int) (CommentResponse, error) {
	comment, err := s.repo.GetById(id)
	if err != nil {
		return CommentResponse{}, err
	}
	return comment, nil
}

func (s *CommentService) Update(id int, req UpdateCommentRequest) error {
	return s.repo.Update(id, req)
}

func (s *CommentService) Delete(id int) error {
	return s.repo.Delete(id)
}
