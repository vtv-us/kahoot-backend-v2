package cloudinary

import (
	"context"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

type CloudinaryService struct {
	Service      *cloudinary.Cloudinary
	UploadFolder string
}

func NewCloudinaryService(c *utils.Config) (CloudinaryService, error) {
	cloudinary, err := cloudinary.NewFromURL(c.CloudinaryUrl)
	if err != nil {
		return CloudinaryService{}, err
	}
	return CloudinaryService{
		Service:      cloudinary,
		UploadFolder: c.CloudinaryUploadFolder,
	}, nil
}

func (s *CloudinaryService) UploadImage(input interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uploadResult, err := s.Service.Upload.Upload(ctx, input, uploader.UploadParams{Folder: s.UploadFolder})
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}
