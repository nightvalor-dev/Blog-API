package cloudinary

import (
	"Project2-v7/config"

	"github.com/cloudinary/cloudinary-go/v2"
)

func NewClient(cfg *config.Config) (*cloudinary.Cloudinary, error) {
	return cloudinary.NewFromParams(
		cfg.CloudinaryCloudName,
		cfg.CloudinaryApiKey,
		cfg.CloudinaryApiSecret,
	)
}
