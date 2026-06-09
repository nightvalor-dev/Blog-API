package media

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type mediaRepository struct {
	db *pgxpool.Pool
}

func NewMediaRepository(
	db *pgxpool.Pool,
) MediaRepository {

	return &mediaRepository{
		db: db,
	}
}

func (r *mediaRepository) Save(blogId int, url string, publicId string) (int, error) {
	var mediaId int
	query := "INSERT INTO media (blog_id, url, public_id) VALUES ($1,$2,$3) RETURNING media_id"
	err := r.db.QueryRow(context.Background(), query, blogId, url, publicId).Scan(&mediaId)
	return mediaId, err
}

func (r *mediaRepository) GetByBlogId(blogId int) ([]Media, error) {
	query := "SELECT media_id, blog_id, url, public_id, created_at FROM media WHERE blog_id = $1"
	rows, err := r.db.Query(context.Background(), query, blogId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Media
	for rows.Next() {
		var m Media
		err := rows.Scan(&m.MediaId, &m.BlogId, &m.Url, &m.PublicId, &m.CreatedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, m)
	}

	return result, nil
}

func (r *mediaRepository) Delete(mediaId int) (string, error) {
	var publicId string
	query := "DELETE FROM media WHERE media_id = $1 RETURNING public_id"
	err := r.db.QueryRow(context.Background(), query, mediaId).Scan(&publicId)
	if err != nil {
		return "", err
	}
	return publicId, nil
}

func (r *mediaRepository) UpdateURLAndPublicID(mediaId int, url, publicID string) error {
	query := "UPDATE media SET url = $1, public_id = $2 WHERE media_id = $3"
	_, err := r.db.Exec(context.Background(), query, url, publicID, mediaId)
	return err
}
