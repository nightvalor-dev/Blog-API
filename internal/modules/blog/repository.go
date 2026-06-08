package blog

import "Project2-v7/pkg/utils"

type BlogRepository interface {
	Create(blog CreateBlogRequest) (int, error)
	GetAllFiltered(filter BlogFilter, pagination utils.PaginationParams) ([]Blog, int, error) // filtered route
	GetById(id int) (Blog, error)
	Update(id int, newBlog UpdateBlogRequest) error
	Delete(id int) error
	AddCategory(blogId int, categoryId int) error
	AddTag(blogId int, tagId int) error
}
