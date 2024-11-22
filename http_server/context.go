package http_server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type HTTPContext struct {
	w http.ResponseWriter
	*http.Request
}

func NewHTTPContext(w http.ResponseWriter, r *http.Request) *HTTPContext {
	return &HTTPContext{w, r}
}

func (c *HTTPContext) sendResponseJSON(statusCode int, data any) {
	c.w.Header().Add("Content-Type", "application/json; charset=utf-8")
	c.w.WriteHeader(statusCode)

	encoder := json.NewEncoder(c.w)
	err := encoder.Encode(data)
	if err != nil {
		slog.Warn(err.Error())
	}
}

func (c *HTTPContext) sendErrorResponse(statusCode int, errCode int, msg string) {
	c.sendResponseJSON(statusCode, ApiResponse{
		Success: false,
		Code:    errCode,
		Data:    msg,
	})
}
