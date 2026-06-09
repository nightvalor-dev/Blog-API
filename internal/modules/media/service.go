package media

import (
	"bytes"
	"context"
	"io"
	"log"
	"mime/multipart"

	sharedcloudinary "Project2-v7/internal/shared/cloudinary"
)

type MediaService struct {
	repo       MediaRepository
	cloudinary *sharedcloudinary.Service
}

func NewMediaService(
	repo MediaRepository,
	cloudinary *sharedcloudinary.Service,
) *MediaService {

	return &MediaService{
		repo:       repo,
		cloudinary: cloudinary,
	}
}

func (s *MediaService) Upload(ctx context.Context, blogId int, file multipart.File) (MediaResponse, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return MediaResponse{}, err
	}

	mediaId, err := s.repo.Save(blogId, "", "")
	if err != nil {
		return MediaResponse{}, err
	}

	go func() {
		reader := bytes.NewReader(fileBytes)
		url, publicID, err := s.cloudinary.Upload(context.Background(), reader)

		if err != nil {
			log.Printf("cloudinary upload failed for media %d: %v", mediaId, err)
			return
		}

		if err := s.repo.UpdateURLAndPublicID(mediaId, url, publicID); err != nil {
			log.Printf("failed to update media %d after upload: %v", mediaId, err)
		}
	}()

	return MediaResponse{
		MediaId: mediaId,
		BlogId:  blogId,
		Url:     "", // client gets empty URL, polls or waits
	}, nil
}

func (s *MediaService) GetByBlogId(ctx context.Context, blogId int) ([]MediaResponse, error) {
	mediaList, err :=
		s.repo.GetByBlogId(blogId)

	if err != nil {
		return nil, err
	}

	result := make([]MediaResponse, len(mediaList))
	for i, m := range mediaList {

		result[i] = MediaResponse{
			MediaId: m.MediaId,
			BlogId:  m.BlogId,
			Url:     m.Url,
		}
	}

	return result, nil
}

func (s *MediaService) Delete(ctx context.Context, mediaId int) error {
	publicId, err := s.repo.Delete(mediaId)
	if err != nil {
		return err
	}

	if publicId == "" {
		return nil
	}

	return s.cloudinary.Delete(ctx, publicId)
}
