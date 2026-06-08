package category

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type catRepository struct {
	db *pgxpool.Pool
}

func NewCatRepository(db *pgxpool.Pool) CategoryRepository {
	return &catRepository{db: db}
}

func (cat *catRepository) FindOrCreate(req CategoryRequest) (int, error) {
	var categoryId int
	err := cat.db.QueryRow(context.Background(),
		`SELECT category_id FROM categories WHERE category_name = $1`,
		req.CategoryName,
	).Scan(&categoryId)

	if err == nil {
		return categoryId, nil // already exists
	}

	err = cat.db.QueryRow(context.Background(),
		`INSERT INTO categories (category_name, description) VALUES ($1, $2) RETURNING category_id`,
		req.CategoryName, req.Description,
	).Scan(&categoryId)

	return categoryId, err
}

func (cat *catRepository) GetAll() ([]CategoryResponse, error) {
	query := "SELECT category_name, description FROM categories"
	rows, err := cat.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []CategoryResponse
	for rows.Next() {
		var entity CategoryResponse
		err = rows.Scan(&entity.CategoryName, &entity.Description)
		if err != nil {
			return nil, err
		}

		result = append(result, entity)
	}

	return result, nil
}

func (cat *catRepository) GetById(id int) (CategoryResponse, error) {
	var entity CategoryResponse
	query := `SELECT category_name, description FROM categories WHERE category_id = $1`
	err := cat.db.QueryRow(context.Background(), query, id).Scan(&entity.CategoryName, &entity.Description)

	if err != nil {
		return CategoryResponse{}, err
	}

	return entity, nil
}

func (cat *catRepository) Update(id int, newCategory CategoryRequest) error {
	args := []any{}
	argIdx := 1
	query := "UPDATE categories SET "
	sep := ""

	if newCategory.CategoryName != "" {
		query += fmt.Sprintf("%s category_name = $%d", sep, argIdx)
		args = append(args, newCategory.CategoryName)
		argIdx++
		sep = ","
	}

	if newCategory.Description != "" {
		query += fmt.Sprintf("%s description = $%d", sep, argIdx)
		args = append(args, newCategory.Description)
		argIdx++
		sep = ","
	}

	if len(args) == 0 {
		return nil
	}

	query += fmt.Sprintf(" WHERE category_id = $%d", argIdx)
	args = append(args, id)

	_, err := cat.db.Exec(context.Background(), query, args...)
	return err
}

func (cat *catRepository) Delete(id int) error {
	query := `DELETE FROM categories WHERE category_id = $1`
	_, err := cat.db.Exec(context.Background(), query, id)
	return err
}

func (c *catRepository) GetByBlogId(blogId int) ([]CategoryResponse, error) {
	query := `
		SELECT c.category_name, c.description
		FROM categories c
		JOIN blog_categories bc ON bc.category_id = c.category_id
		WHERE bc.blog_id = $1`

	rows, err := c.db.Query(context.Background(), query, blogId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []CategoryResponse
	for rows.Next() {
		var cat CategoryResponse
		err = rows.Scan(&cat.CategoryName, &cat.Description)
		if err != nil {
			return nil, err
		}
		result = append(result, cat)
	}
	return result, nil
}
