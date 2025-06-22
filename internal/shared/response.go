package shared

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpError struct {
	Code    int
	Message string
}

type HttpRes struct {
	Code    int
	Payload any
}

func (e HttpError) Error() string {
	return e.Message
}

func ResHttpError(c *gin.Context, httpError *HttpError) {
	c.JSON(httpError.Code, map[string]string{"error": httpError.Message})
}

func ResTryHttpError(c *gin.Context, err error) {
	httpError, ok := err.(*HttpError)
	if !ok {
		ResInternalServerErrorDef(c)
		return
	}

	c.JSON(httpError.Code, map[string]string{"error": httpError.Message})
}

func ResOK(c *gin.Context, payload any) {
	c.JSON(http.StatusOK, payload)
}

func ResCreated(c *gin.Context, payload any) {
	c.JSON(http.StatusCreated, payload)
}

func ResAccepted(c *gin.Context, payload any) {
	c.JSON(http.StatusAccepted, payload)
}

func ResNonAuthoritativeInfo(c *gin.Context, payload any) {
	c.JSON(http.StatusNonAuthoritativeInfo, payload)
}

func ResNoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

func ResBadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, map[string]string{"error": msg})
}

func ResUnauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, map[string]string{"error": msg})
}

func ResForbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, map[string]string{"error": msg})
}

func ResNotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, map[string]string{"error": msg})
}

func ResMethodNotAllowed(c *gin.Context, msg string) {
	c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": msg})
}

func ResNotAcceptable(c *gin.Context, msg string) {
	c.JSON(http.StatusNotAcceptable, map[string]string{"error": msg})
}

func ResConflict(c *gin.Context, msg string) {
	c.JSON(http.StatusConflict, map[string]string{"error": msg})
}

func ResInternalServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, map[string]string{"error": msg})
}

func ResInternalServerErrorDef(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, map[string]string{"error": ErrInternal})
}

func ResNotImplemented(c *gin.Context, msg string) {
	c.JSON(http.StatusNotImplemented, map[string]string{"error": msg})
}

func OK(payload any) *HttpRes {
	return &HttpRes{
		Code:    http.StatusOK,
		Payload: payload,
	}
}

func Created(payload any) *HttpRes {
	return &HttpRes{
		Code:    http.StatusCreated,
		Payload: payload,
	}
}

func Accepted(payload any) *HttpRes {
	return &HttpRes{
		Code:    http.StatusAccepted,
		Payload: payload,
	}
}

func NonAuthoritativeInfo(payload any) *HttpRes {
	return &HttpRes{
		Code:    http.StatusNonAuthoritativeInfo,
		Payload: payload,
	}
}

func NoContent() *HttpRes {
	return &HttpRes{
		Code:    http.StatusNoContent,
		Payload: nil,
	}
}

func BadRequest(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusBadRequest,
	}
}

func Unauthorized(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusUnauthorized,
	}
}

func Forbidden(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusForbidden,
	}
}

func NotFound(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusNotFound,
	}
}

func MethodNotAllowed(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusMethodNotAllowed,
	}
}

func NotAcceptable(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusNotAcceptable,
	}
}

func Conflict(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusConflict,
	}
}

func InternalServerError(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusInternalServerError,
	}
}

func InternalServerErrorDef() *HttpError {
	return &HttpError{
		Message: ErrInternal,
		Code:    http.StatusInternalServerError,
	}
}

func NotImplemented(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusNotImplemented,
	}
}
