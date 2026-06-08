package blog

type BlogStatus string

const (
	StatusDraft     BlogStatus = "drafted"
	StatusPublished BlogStatus = "published"
	StatusArchived  BlogStatus = "archived"
)
