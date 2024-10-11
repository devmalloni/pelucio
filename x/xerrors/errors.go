package xerrors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ApplicationError struct {
	HTTPStatus  int    `json:"-"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
	Data        any    `json:"data,omitempty"`
	InnerErr    error  `json:"innerError,omitempty"`
}

func (p ApplicationError) Error() string {
	return p.Description
}

func (p ApplicationError) Unwrap() error {
	return p.InnerErr
}

func New(code string) ApplicationError {
	return ApplicationError{
		Code:       code,
		HTTPStatus: http.StatusBadRequest, // by default, we consider every error as bad request
	}
}

func (p ApplicationError) WithDescription(desc string, args ...any) ApplicationError {
	p.Description = fmt.Sprintf(desc, args...)
	return p
}

func (p ApplicationError) WithData(dt any) ApplicationError {
	p.Data = dt
	return p
}

func (p ApplicationError) WithHttpStatus(status int) ApplicationError {
	p.HTTPStatus = status
	return p
}

func (p ApplicationError) WithError(err error) ApplicationError {
	p.InnerErr = err
	return p
}

func (p ApplicationError) Wrap(description string, args ...any) ApplicationError {
	return New(p.Code).WithDescription(description, args...).WithError(p)
}

func TryGetErrorFromResponse(res *http.Response) (*ApplicationError, []byte, bool) {
	responsebody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, false
	}
	var applicationError ApplicationError
	err = json.Unmarshal(responsebody, &applicationError)
	if err != nil {
		return nil, nil, false
	}

	return &applicationError, responsebody, true
}
