package category

type CategoryRepository interface {
	FindOrCreate(category CategoryRequest) (int, error)
	GetAll() ([]CategoryResponse, error)
	GetById(id int) (CategoryResponse, error)
	Update(id int, newCategory CategoryRequest) error
	Delete(id int) error
	GetByBlogId(id int) ([]CategoryResponse, error)
}
