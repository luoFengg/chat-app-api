package upload

import "github.com/gin-gonic/gin"

// UploadController handles file upload HTTP requests
type UploadController interface {
	// UploadFile handles POST /api/v1/upload
	UploadFile(ctx *gin.Context)
}