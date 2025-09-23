// utils/fileupload.go
package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/wildanasyrof/backend-topup/pkg/storage"
)

const MaxImageSize = 2 * 1024 * 1024 // 2MB

var allowedExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

// UploadImage validates then saves an image using storage, returns public URL.
func UploadImage(
	file *multipart.FileHeader,
	st storage.LocalStorage, // contoh: "/uploads"
) (string, error) {
	if file.Size > MaxImageSize {
		return "", fmt.Errorf("file too large (max 2MB)")
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExt[ext] {
		return "", fmt.Errorf("unsupported image type (jpg/jpeg/png/webp only)")
	}

	filename, err := st.Save(file)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}
	return storage.PublicURL("/uploads", filename), nil
}
