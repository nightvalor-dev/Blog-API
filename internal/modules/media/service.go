package media

import (
	"context"
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
	url, publicID, err :=
		s.cloudinary.Upload(ctx, file)

	if err != nil {
		return MediaResponse{}, err
	}

	mediaId, err := s.repo.Save(blogId, url, publicID)
	if err != nil {
		return MediaResponse{}, err
	}

	return MediaResponse{
		MediaId: mediaId,
		BlogId:  blogId,
		Url:     url,
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
