package cloudinary

import (
	"context"

	cld "github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Service struct {
	client *cld.Cloudinary
}

func NewService(client *cld.Cloudinary) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) Upload(ctx context.Context, file interface{}) (string, string, error) {
	result, err := s.client.Upload.Upload(ctx, file, uploader.UploadParams{Folder: "blog-media"})
	if err != nil {
		return "", "", err
	}

	return result.SecureURL, result.PublicID, nil
}

func (s *Service) Delete(ctx context.Context, publicId string) error {
	_, err := s.client.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicId})
	return err
}
