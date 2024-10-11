package xerrors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HandlerFuncWithError func(c *gin.Context) error

// HandleWithError is a helper function that wraps an endpoint with error return and handles every error that can occur in the request.
func HandleWithError(h HandlerFuncWithError) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := h(c)
		GinError(err, c)
	}
}

// GinError is a helper function that can be used to check if there's an error e.g if GinError(err, c) { return }
func GinError(err error, c *gin.Context) bool {
	if err == nil {
		return false
	}

	c.Set("application-error", err)

	return true
}

// HandleErrorMiddleware is a middleware that takes the error
func HandleErrorMiddleware(c *gin.Context) {
	c.Next()
	// check if there is some error on context
	ginError, exists := c.Get("application-error")
	if !exists {
		return
	}

	// Cast it to error
	err := ginError.(error)
	// and ensure that is an app error. This will cast a common error to app error
	appError := Ensure(err)

	log.Error().Any("err", appError).Msg("error on handling")
	// write the error on json body
	c.JSON(appError.HTTPStatus, appError)
}

// Ensure returns the error as ApplicationError or
// wraps it on a new generic ApplicationError
func Ensure(err error) (appErr ApplicationError) {
	if errors.As(err, &appErr) {
		// if status is not set, then we can
		// assume that the error is a bad request
		//
		// In this kind of situation, we are assuming
		// that errors that isn't handled by our application
		// is actually internal server errors.
		if appErr.HTTPStatus == 0 {
			appErr = appErr.WithHttpStatus(http.StatusBadRequest)
		}
		return
	}

	appErr = New("ErrInternalServerError").WithDescription(err.Error()).WithHttpStatus(http.StatusInternalServerError)

	return
}
