package blog

import (
	"Project2-v7/pkg/utils"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type blogRepository struct {
	db *pgxpool.Pool
}

func NewBlogRepository(db *pgxpool.Pool) BlogRepository {
	return &blogRepository{db: db}
}

func (b *blogRepository) Create(blog CreateBlogRequest) (int, error) {
	var blogId int
	query := "INSERT INTO blogs (title, content, status, user_id) VALUES ($1, $2, $3, $4) RETURNING blog_id"
	err := b.db.QueryRow(context.Background(), query, blog.Title, blog.Content, blog.Status, blog.UserId).Scan(&blogId)
	if err != nil {
		return 0, err
	}
	return blogId, nil
}

func (b *blogRepository) GetAllFiltered(filter BlogFilter, pagination utils.PaginationParams) ([]Blog, int, error) {
	args := []any{}
	argIdx := 1
	conditions := []string{}

	// --- filtering ---
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}

	if filter.UserId > 0 {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIdx))
		args = append(args, filter.UserId)
		argIdx++
	}

	if filter.CategoryName != "" {
		conditions = append(conditions, fmt.Sprintf(
			`EXISTS (
                SELECT 1 FROM blog_categories bc
                JOIN categories c ON c.category_id = bc.category_id
                WHERE bc.blog_id = blogs.blog_id
                AND LOWER(c.category_name) = LOWER($%d)
            )`, argIdx))
		args = append(args, filter.CategoryName)
		argIdx++
	}

	if filter.TagName != "" {
		conditions = append(conditions, fmt.Sprintf(
			`EXISTS (
                SELECT 1 FROM blog_tags bt
                JOIN tags t ON t.tag_id = bt.tag_id
                WHERE bt.blog_id = blogs.blog_id
                AND LOWER(t.tag_name) = LOWER($%d)
            )`, argIdx))
		args = append(args, filter.TagName)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// --- sorting ---
	allowedSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"title":      true,
	}

	sortBy := "created_at"
	if allowedSortFields[filter.SortBy] {
		sortBy = filter.SortBy
	}

	order := "desc"
	if strings.ToLower(filter.Order) == "asc" {
		order = "asc"
	}

	// --- total count ---
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM blogs %s", whereClause)
	var total int
	if err := b.db.QueryRow(context.Background(), countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// --- paginated data ---
	dataQuery := fmt.Sprintf(
		`SELECT blog_id, title, content, status, user_id FROM blogs %s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		whereClause, sortBy, order, argIdx, argIdx+1,
	)
	args = append(args, pagination.Limit, pagination.Offset)

	rows, err := b.db.Query(context.Background(), dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []Blog
	for rows.Next() {
		var entity Blog
		if err = rows.Scan(
			&entity.BlogId, &entity.Title, &entity.Content,
			&entity.Status, &entity.UserId,
		); err != nil {
			return nil, 0, err
		}
		result = append(result, entity)
	}

	return result, total, nil
}

func (b *blogRepository) GetById(id int) (Blog, error) {
	var entity Blog
	query := `SELECT blog_id, title, content, status, user_id FROM blogs WHERE blog_id = $1`
	err := b.db.QueryRow(context.Background(), query, id).Scan(
		&entity.BlogId, &entity.Title, &entity.Content, &entity.Status, &entity.UserId,
	)

	if err != nil {
		return Blog{}, err
	}
	return entity, nil
}

func (b *blogRepository) Update(id int, newBlog UpdateBlogRequest) error {
	args := []any{}
	argIdx := 1
	query := "UPDATE blogs SET "
	sep := ""

	if newBlog.Title != "" {
		query += fmt.Sprintf("%s title = $%d", sep, argIdx)
		args = append(args, newBlog.Title)
		argIdx++
		sep = ","
	}

	if newBlog.Content != "" {
		query += fmt.Sprintf("%s content = $%d", sep, argIdx)
		args = append(args, newBlog.Content)
		argIdx++
		sep = ","
	}

	if newBlog.Status != "" {
		query += fmt.Sprintf("%s status = $%d", sep, argIdx)
		args = append(args, newBlog.Status)
		argIdx++
		sep = ","
	}

	if len(args) == 0 {
		return nil
	}

	query += fmt.Sprintf(" WHERE blog_id = $%d", argIdx)
	args = append(args, id)
	_, err := b.db.Exec(context.Background(), query, args...)
	return err
}

func (b *blogRepository) Delete(id int) error {
	query := `DELETE FROM blogs WHERE blog_id = $1`
	_, err := b.db.Exec(context.Background(), query, id)
	return err
}

func (b *blogRepository) AddCategory(blogId int, categoryId int) error {
	query := `INSERT INTO blog_categories (blog_id, category_id) VALUES ($1, $2)`
	_, err := b.db.Exec(context.Background(), query, blogId, categoryId)
	return err
}

func (b *blogRepository) AddTag(blogId int, tagId int) error {
	query := `INSERT INTO blog_tags (blog_id, tag_id) VALUES ($1, $2)`
	_, err := b.db.Exec(context.Background(), query, blogId, tagId)
	return err
}
