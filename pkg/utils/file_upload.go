// utils/fileupload.go
package utils

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
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
		return "", apperror.New(apperror.CodeBadRequest, "image too large (max 2mb)", errors.New("file too large"))
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExt[ext] {
		return "", apperror.New(apperror.CodeBadRequest, "unsupported image type (jpg/jpeg/png/webp only)", errors.New("file error"))
	}

	filename, err := st.Save(file)
	if err != nil {
		return "", apperror.New(apperror.CodeInternal, "failed to save file: %w", errors.New("upload error"))
	}
	return storage.PublicURL("/uploads", filename), nil
}
