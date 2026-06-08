package redis

import "fmt"

const (
	BlogList     = "blogs:list"
	CategoryList = "categories:list"
	TagList      = "tags:list"
)

func BlogKey(id int) string     { return fmt.Sprintf("blogs:%d", id) }
func CategoryKey(id int) string { return fmt.Sprintf("categories:%d", id) }
func TagKey(id int) string      { return fmt.Sprintf("tags:%d", id) }

// BlogListKey builds a cache key only for unfiltered, default-sorted requests.
// Keyed only on page and limit since there are no filters.
func BlogListKey(page, limit int) string {
	return fmt.Sprintf("blogs:list:page=%d:limit=%d", page, limit)
}
