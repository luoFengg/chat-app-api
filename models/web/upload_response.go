package web

// UploadResult represents the result of a file upload
type UploadResult struct {
    URL      string `json:"url"`
    Filename string `json:"filename"`
    Size     int64  `json:"size"`
    MimeType string `json:"mime_type"`
}
