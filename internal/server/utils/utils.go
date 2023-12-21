package utils

import (
	"encoding/json"
	"net/http"

	"github.com/kmg7/fson/internal/err"
	"github.com/kmg7/fson/internal/logger"
)

const (
	CodeUnexpectedError  = "unexpected-internal"
	CodeInvalidBodyError = "invalid-body"
	HeaderContentType    = "Content-Type"
	ContentTypeJson      = "application/json"
	ContentTypeTextPlain = "text/plain"
)

func ErrorResponse(w http.ResponseWriter, r *http.Request, code int, err *err.AppError) {
	b, _ := json.Marshal(ErrorResponseBody{Code: err.Code, Err: err.Messages})
	if err.Internal {
		logger.Error(err)
		w.WriteHeader(500)
	} else {
		w.WriteHeader(code)
	}
	SetContentType(w, ContentTypeJson)
	w.Write(b)
}

func UnexpectedInternalErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	b, _ := json.Marshal(ErrorResponseBody{Code: CodeUnexpectedError, Err: []string{}})
	w.WriteHeader(500)
	logger.Error(err.Error())
	SetContentType(w, ContentTypeJson)
	w.Write(b)
}

func InvalidBodyErrorResponse(w http.ResponseWriter, r *http.Request, messages []string) {
	b, _ := json.Marshal(ErrorResponseBody{Code: CodeInvalidBodyError, Err: messages})
	w.WriteHeader(400)
	SetContentType(w, ContentTypeJson)
	w.Write(b)
}

func SetContentType(w http.ResponseWriter, ct string) {
	w.Header().Set(HeaderContentType, ct)
}

type ErrorResponseBody struct {
	Code string   `json:"code"`
	Err  []string `json:"errors"`
}
