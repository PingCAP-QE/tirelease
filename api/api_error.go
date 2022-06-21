package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// error return struct definition
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err APIError) Error() string {
	return err.Message
}

// middleware error handler in server package
func APIErrorJSONReporter() gin.HandlerFunc {
	return APIErrorJSONReporterHandler(gin.ErrorTypeAny)
}

func APIErrorJSONReporterHandler(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		innerErrors := c.Errors.ByType(errType)

		if len(innerErrors) > 0 {
			err := innerErrors[0].Err
			parsedError := &APIError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			log.Println(parsedError)

			// var parsedError *APIError
			// switch err.(type) {
			// case *APIError:
			// 	parsedError = err.(*APIError)
			// default:
			// 	parsedError = &APIError{
			// 		Code:    http.StatusInternalServerError,
			// 		Message: "Internal Server Error",
			// 	}
			// }

			// Put the error into response
			c.IndentedJSON(parsedError.Code, parsedError)
			c.Abort()
			return
		}

	}
}
