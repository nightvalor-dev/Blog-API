package category

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(req CategoryRequest) (int, error) {
	return s.repo.FindOrCreate(req)
}

func (s *CategoryService) GetAll() ([]CategoryResponse, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) GetById(id int) (CategoryResponse, error) {
	cat, err := s.repo.GetById(id)
	if err != nil {
		return CategoryResponse{}, err
	}
	return cat, nil
}

func (s *CategoryService) Update(id int, req CategoryRequest) error {
	return s.repo.Update(id, req)
}

func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
