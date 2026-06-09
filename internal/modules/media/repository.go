package media

type MediaRepository interface {
	Save(blogId int, url string, publicId string) (int, error)
	GetByBlogId(blogId int) ([]Media, error)
	Delete(mediaId int) (string, error)
	UpdateURLAndPublicID(mediaId int, url, publicID string) error
}
