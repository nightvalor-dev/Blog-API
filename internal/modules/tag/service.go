package tag

type TagService struct {
	repo TagRepository
}

func NewTagService(repo TagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) Create(req TagRequest) (int, error) {
	return s.repo.FindOrCreate(req)
}

func (s *TagService) GetAll() ([]TagResponse, error) {
	tags, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	var result []TagResponse
	for _, ta := range tags {
		result = append(result, ta)
	}
	return result, nil
}

func (s *TagService) GetById(id int) (TagResponse, error) {
	ta, err := s.repo.GetById(id)
	if err != nil {
		return TagResponse{}, err
	}
	return ta, nil
}

func (s *TagService) Update(id int, req TagRequest) error {
	return s.repo.Update(id, req)
}

func (s *TagService) Delete(id int) error {
	return s.repo.Delete(id)
}
