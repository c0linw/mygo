package http_server

import (
	"encoding/json"
	"errors"
	"mygo/errval"
	"mygo/validation"
	"net/http"
)

type SimpleWrapFunc func(r *HTTPContext) (any, error)
type InterfaceWrapFunc[T any] func(val *T, r *HTTPContext) (any, error)
type ApiResponse struct {
	Success bool `json:"success"`
	Code    int  `json:"code"`
	Data    any  `json:"data"`
}

func SimpleWrapper(apiFunc SimpleWrapFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewHTTPContext(w, r)
		res, err := apiFunc(ctx)
		if err != nil {
			var apiError *errval.ApiError
			if errors.As(err, &apiError) {
				ctx.sendErrorResponse(http.StatusBadRequest, apiError.Code, apiError.Error())
			}
			ctx.sendErrorResponse(http.StatusBadRequest, http.StatusBadRequest, err.Error())
			return
		}
		ctx.sendResponseJSON(http.StatusOK, ApiResponse{Success: true, Code: http.StatusOK, Data: res})
	}
}

func JSONRequestBodyWrapper[T any](apiFunc InterfaceWrapFunc[T]) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewHTTPContext(w, r)
		if r.ContentLength <= 0 {
			ctx.sendErrorResponse(http.StatusBadRequest, http.StatusBadRequest, "missing request body")
			return
		}
		val, err := decodeBody[T](r)
		if err != nil {
			ctx.sendErrorResponse(http.StatusBadRequest, http.StatusBadRequest, err.Error())
			return
		}
		res, err := apiFunc(val, ctx)
		if err != nil {
			var apiError *errval.ApiError
			if errors.As(err, &apiError) {
				ctx.sendErrorResponse(http.StatusBadRequest, apiError.Code, apiError.Error())
			}
			ctx.sendErrorResponse(http.StatusBadRequest, -1, err.Error())
			return
		}
		ctx.sendResponseJSON(http.StatusOK, ApiResponse{Success: true, Data: res})
	}
}

func decodeBody[T any](r *http.Request) (*T, error) {
	body := new(T)
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		return nil, err
	}
	if err := validation.ValidateStruct(body); err != nil {
		return nil, err
	}
	return body, nil
}
