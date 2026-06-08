package tag

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type tagRepository struct {
	db *pgxpool.Pool
}

func NewTagRepository(db *pgxpool.Pool) TagRepository {
	return &tagRepository{db: db}
}

func (tr *tagRepository) FindOrCreate(req TagRequest) (int, error) {
	var tagId int
	err := tr.db.QueryRow(context.Background(),
		`SELECT tag_id FROM tags WHERE tag_name = $1`,
		req.TagName,
	).Scan(&tagId)

	if err == nil {
		return tagId, nil // already exists
	}

	err = tr.db.QueryRow(context.Background(), `INSERT INTO tags (tag_name) VALUES ($1) RETURNING tag_id`,
		req.TagName).Scan(&tagId)

	return tagId, err
}

func (tr *tagRepository) GetAll() ([]TagResponse, error) {
	query := `SELECT tag_name FROM tags`
	rows, err := tr.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TagResponse
	for rows.Next() {
		var entity TagResponse
		err = rows.Scan(&entity.TagName)
		if err != nil {
			return nil, err
		}

		result = append(result, entity)
	}
	return result, nil
}

func (tr *tagRepository) GetById(id int) (TagResponse, error) {
	var entity TagResponse
	query := `SELECT tag_name FROM tags WHERE tag_id = $1`
	err := tr.db.QueryRow(context.Background(), query, id).Scan(&entity.TagName)
	if err != nil {
		return TagResponse{}, err
	}
	return entity, nil
}

func (tr *tagRepository) Update(id int, newTag TagRequest) error {
	args := []any{}
	argIdx := 1
	query := "UPDATE tags SET "
	sep := ""

	if newTag.TagName != "" {
		query += fmt.Sprintf("%s tag_name = $%d", sep, argIdx)
		args = append(args, newTag.TagName)
		argIdx++
		sep = ","
	}

	if len(args) == 0 {
		return nil
	}

	query += fmt.Sprintf(" WHERE tag_id = $%d", argIdx)
	args = append(args, id)

	_, err := tr.db.Exec(context.Background(), query, args...)
	return err
}

func (tr *tagRepository) Delete(id int) error {
	query := `DELETE FROM tags WHERE tag_id = $1`
	_, err := tr.db.Exec(context.Background(), query, id)
	return err
}

func (t *tagRepository) GetByBlogId(blogId int) ([]TagResponse, error) {
	query := `
		SELECT ta.tag_name
		FROM tags ta
		JOIN blog_tags bt ON bt.tag_id = ta.tag_id
		WHERE bt.blog_id = $1`

	rows, err := t.db.Query(context.Background(), query, blogId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TagResponse
	for rows.Next() {
		var tag TagResponse
		err = rows.Scan(&tag.TagName)
		if err != nil {
			return nil, err
		}
		result = append(result, tag)
	}
	return result, nil
}
