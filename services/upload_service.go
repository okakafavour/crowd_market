package services

import (
	"CROWD_MARKET/config"
	"context"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// ✅ Upload image to Cloudinary
func UploadToCloudinary(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	uploadResult, err := config.Cloud.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: fileHeader.Filename,
		Folder:   "crowd_market/products",
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

// 🧹 Delete an image from Cloudinary using its public ID
func DeleteFromCloudinary(imageURL string) error {
	ctx := context.Background()

	// Extract the public ID from the image URL
	// Example: https://res.cloudinary.com/.../crowd_market/products/image.jpg
	parts := strings.Split(imageURL, "/")
	if len(parts) == 0 {
		return nil
	}

	publicIDWithExt := parts[len(parts)-1]
	publicID := strings.TrimSuffix(publicIDWithExt, filepath.Ext(publicIDWithExt))

	// Prepend folder name
	fullPublicID := "crowd_market/products/" + publicID

	_, err := config.Cloud.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: fullPublicID,
	})
	return err
}
