package response

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	apperr "github.com/farismfirdaus/plant-nursery/errors"
)

// ResponseSuccess contains response success body
type ResponseSuccess struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

// ResponseError contains response error body
type ResponseError struct {
	Error []Error `json:"error"`
}

// Error contains error body
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var (
	// InternalServerError internal server error
	InternalServerError = ResponseError{Error: []Error{{"internal server error", 10000}}}
)

// BuildSuccess build a success response
func BuildSuccess(c *gin.Context, code int, data any) {
	if reflect.TypeOf(data).Kind() == reflect.String {
		// return response message
		c.JSON(code, ResponseSuccess{Message: fmt.Sprint(data)})
		return
	}
	// return response data
	c.JSON(code, ResponseSuccess{Data: data})
}

// BuildErrors build a error response
// Note:
//   - all of error's code should be the same, else returning internal server error
//   - specified error on errors package
func BuildErrors(c *gin.Context, errs ...error) {
	respErrors := ResponseError{}

	// init var to check error's code similarity
	var (
		respStatusCode    int
		respStatusCodeSet bool
	)

	// if no error passed, return internal server error
	if len(errs) <= 0 {
		respStatusCode = http.StatusInternalServerError
		respErrors = InternalServerError
	}

	for _, err := range errs {
		var (
			model = apperr.Error{}
			e     = Error{}
		)

		// check if error implement Error struct
		if errors.As(err, &model) {
			e = Error{Message: model.Message, Code: model.Code}
		} else {
			respStatusCode = http.StatusInternalServerError
			respErrors = InternalServerError
			break
		}

		// check code similarity
		if respStatusCodeSet && respStatusCode != model.StatusCode {
			respStatusCode = http.StatusInternalServerError
			respErrors = InternalServerError
			break
		} else {
			respStatusCodeSet = true
			respStatusCode = model.StatusCode
		}

		respErrors.Error = append(respErrors.Error, e)
		c.Error(err)
	}

	c.JSON(respStatusCode, respErrors)
}
