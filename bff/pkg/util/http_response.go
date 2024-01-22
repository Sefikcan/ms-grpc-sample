package util

import (
	"errors"
	"fmt"
)

const (
	ErrBadRequest = "Bad request"
	ErrNotFound   = "Not found"
)

var (
	BadRequest          = errors.New("Bad Request")
	NotFound            = errors.New("Not Found")
	InternalServerError = errors.New("Internal Server Error")
	RequestTimeoutError = errors.New("Request Timeout")
)

type HttpResponse interface {
	Status() int
	Error() string
	Causes() interface{}
}

type httpResponse struct {
	ErrStatus int         `json:"status,omitempty"`
	ErrError  string      `json:"error,omitempty"`
	ErrCauses interface{} `json:"-"`
}

func (h httpResponse) Status() int {
	return h.ErrStatus
}

func (h httpResponse) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", h.ErrStatus, h.ErrError, h.ErrCauses)
}

func (h httpResponse) Causes() interface{} {
	return h.ErrCauses
}

func NewHttpResponse(status int, err string, causes interface{}) HttpResponse {
	return &httpResponse{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}
