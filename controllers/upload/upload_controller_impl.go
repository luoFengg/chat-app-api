package upload

import (
	"chatapp-api/models/web"
	uploadService "chatapp-api/services/upload"
	"net/http"

	"github.com/gin-gonic/gin"
)

// uploadControllerImpl implements UploadController interface
type uploadControllerImpl struct {
	uploadService uploadService.UploadService
}

// NewUploadController create new instance of UploadController
func NewUploadController(uploadService uploadService.UploadService) UploadController {
	return &uploadControllerImpl{uploadService: uploadService}
}

// UploadFile handles POST /api/v1/upload
func (controller *uploadControllerImpl) UploadFile(ctx *gin.Context) {
	// 1. Get file from multipart form
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrorResponse{
			Success: false,
			Message: "File is required. Use 'file' as the form field name",
			Error: err.Error(),
		})
		return
	}
	defer file.Close()

	// 2. Call upload service
	result, err := controller.uploadService.UploadFile(ctx.Request.Context(), file, header)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 3. Return success response
	ctx.JSON(http.StatusOK, web.ApiResponse{
		Success: true,
		Message: "File uploaded successfully",
		Data: result,
	})
}