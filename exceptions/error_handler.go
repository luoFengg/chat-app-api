package exceptions

import (
	"chatapp-api/models/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a Middleware to Handle Global Error
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Run the Handler First

		// Check if there is an error
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last().Err

			// Handle Based on the Error Type
			switch e := err.(type) {
			case NotFoundError:
				c.JSON(http.StatusNotFound, web.ApiResponse{
					Success: false,
					Message: e.Message,
				})
			case BadRequestError:
				c.JSON(http.StatusBadRequest, web.ApiResponse{
					Success: false,
					Message: e.Message,
				})
			case UnauthorizedError:
				c.JSON(http.StatusUnauthorized, web.ApiResponse{
					Success: false,
					Message: e.Message,
				})
			case ConflictError:
				c.JSON(http.StatusConflict, web.ApiResponse{
					Success: false,
					Message: e.Message,
				})
			case ForbiddenError:
				c.JSON(http.StatusForbidden, web.ApiResponse{
					Success: false,
					Message: e.Message,
				})
			default:
				c.JSON(http.StatusInternalServerError, web.ApiResponse{
					Success: false,
					Message: "Internal server error",
				})
			}
		}
	}
}