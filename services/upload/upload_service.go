package upload

import (
	"chatapp-api/models/web"
	"context"
	"mime/multipart"
)

// UploadService interface for file upload operations
type UploadService interface {
	// UploadFile uploads a file to Supabase Storage and returns the public URL
	UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*web.UploadResult, error)
}