package category

type CategoryRequest struct {
	CategoryName string `json:"category_name"`
	Description  string `json:"description"`
}

type CategoryResponse struct {
	CategoryName string `json:"category_name"`
	Description  string `json:"description"`
}
