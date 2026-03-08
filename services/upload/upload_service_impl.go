package upload

import (
	"chatapp-api/config"
	"chatapp-api/exceptions"
	"chatapp-api/models/web"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/oklog/ulid/v2"
)

// Allowed MIME types
var allowedMimeTypes = map[string]bool{
    // Images
    "image/jpeg": true,
    "image/png":  true,
    "image/gif":  true,
    "image/webp": true,
    // Documents
    "application/pdf":   true,
    "application/msword": true,
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
    // Video
    "video/mp4":  true,
    "video/webm": true,
    // Audio
    "audio/mpeg": true,
    "audio/ogg":  true,
    "audio/wav":  true,
    "audio/webm": true,
}

const maxFileSize = 20 * 1024 * 1024

// uploadServiceImpl implements the UploadService interface
type uploadServiceImpl struct {
	config *config.Config
}

// NewUploadService creates a new instance of UploadService
func NewUploadService(config *config.Config) UploadService {
	return &uploadServiceImpl{config: config}
}


// UploadFile uploads a file to Supabase Storage and returns the public URL
func (service *uploadServiceImpl) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*web.UploadResult, error) {
	// 1. Validate file size
	if header.Size > maxFileSize {
        return nil, exceptions.NewBadRequestError(
            fmt.Sprintf("File too large. Maximum size is %d MB", maxFileSize/(1024*1024)))
    }

	// 2. Detect MIME type from file content
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return nil, exceptions.NewInternalServerError("Failed to read file")
	}
	mimetype := http.DetectContentType(buffer)
	file.Seek(0, io.SeekStart)

	// 3. Validate MIME type
	if !allowedMimeTypes[mimetype] {
		return nil, exceptions.NewBadRequestError(
			fmt.Sprintf("File type '%s' is not allowed", mimetype))
	}
	
	// 4. Generate unique filename
	extension := filepath.Ext(header.Filename)
	if extension == "" {
		extension = getExtensionFromMime(mimetype)
	}
	uniqueFileName := ulid.Make().String() + extension

	// 5. Subfolder based on MIME type
	subfolder := getSubfolder(mimetype)
	objectPath := subfolder + "/" + uniqueFileName

	// 6. Upload to Supabase Storage via HTTP POST
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", service.config.Supabase.URL, service.config.Supabase.Bucket, objectPath)

	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, file)
	if err != nil {
		return nil, exceptions.NewInternalServerError("Failed to create upload request")
	}
	req.Header.Set("Authorization", "Bearer "+service.config.Supabase.Key)
	req.Header.Set("Content-Type", mimetype)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, exceptions.NewInternalServerError("Failed to upload file to storage")
	}
	defer response.Body.Close()

	// 7. Check the response
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(response.Body)
		return nil, exceptions.NewInternalServerError(
			fmt.Sprintf("Storage upload failed (status %d): %s", response.StatusCode, string(body)))
	}

	// 8. Return public URL
	   publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", service.config.Supabase.URL, service.config.Supabase.Bucket, objectPath)

	   return &web.UploadResult{
		URL: publicURL,
		Filename: header.Filename,
		Size: header.Size,
		MimeType: mimetype,
	   }, nil
}

// getSubfoler helper method to determine folder based on MIME prefix
func getSubfolder(mimeType string) string {
    if strings.HasPrefix(mimeType, "image/") {
        return "images"
    }
    if strings.HasPrefix(mimeType, "video/") {
        return "videos"
    }
    if strings.HasPrefix(mimeType, "audio/") {
        return "audio"
    }
    return "files"
}

// getExtensionFromMime helper method to fallback if the file has no extension
func getExtensionFromMime(mimeType string) string {
    extensions := map[string]string{
        "image/jpeg": ".jpg",
        "image/png":  ".png",
        "image/gif":  ".gif",
        "image/webp": ".webp",
        "video/mp4":  ".mp4",
        "video/webm": ".webm",
        "audio/mpeg": ".mp3",
        "audio/ogg":  ".ogg",
        "audio/wav":  ".wav",
        "audio/webm": ".webm",
        "application/pdf": ".pdf",
    }
    if ext, ok := extensions[mimeType]; ok {
        return ext
    }
    return ""
}
