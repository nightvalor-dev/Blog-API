package tag

type TagRepository interface {
	FindOrCreate(tag TagRequest) (int, error)
	GetAll() ([]TagResponse, error)
	GetById(id int) (TagResponse, error)
	Update(id int, newTag TagRequest) error
	Delete(id int) error
	GetByBlogId(id int) ([]TagResponse, error)
}
