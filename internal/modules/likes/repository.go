package likes

type LikeRepository interface {
	Add(blogId, userId int) error
	Remove(blogId, userId int) error
	Count(blogId int) (int, error)
	Exists(blogId, userId int) (bool, error)
}
