package blog

type CategoryDTOLocal struct {
	CategoryName string `json:"category_name"`
	Description  string `json:"description"`
}

type TagDTOLocal struct {
	TagName string `json:"tag_name"`
}

type MediaDTOLocal struct {
	Url string `json:"url"`
}

// POST Operation
type CreateBlogRequest struct {
	Title    string             `json:"title"`
	Content  string             `json:"content"`
	Status   BlogStatus         `json:"status"`
	UserId   int                `json:"user_id"`
	Category []CategoryDTOLocal `json:"category"`
	Tags     []TagDTOLocal      `json:"tags"`
}

// PUT Operation
type UpdateBlogRequest struct {
	Title   string     `json:"title"`
	Content string     `json:"content"`
	Status  BlogStatus `json:"status"`
}

// GET Operation
type BlogResponse struct {
	Title    string             `json:"title"`
	Content  string             `json:"content"`
	Status   BlogStatus         `json:"status"`
	Category []CategoryDTOLocal `json:"category"`
	Tags     []TagDTOLocal      `json:"tags"`
	Images   []MediaDTOLocal    `json:"images"`
}

type BlogFilter struct {
	Status       string
	UserId       int
	SortBy       string
	Order        string
	CategoryName string
	TagName      string
}

type BlogListResponse struct {
	Data       []BlogResponse `json:"data"`
	TotalCount int            `json:"total_count"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
}
