package http_server

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/schema"
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

var decoder = schema.NewDecoder()

func SimpleWrapper(apiFunc SimpleWrapFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewHTTPContext(w, r)
		res, err := apiFunc(ctx)
		if err != nil {
			var apiError *errval.ApiError
			if errors.As(err, &apiError) {
				ctx.sendErrorResponse(http.StatusBadRequest, apiError.Code, apiError.Error())
				return
			}
			ctx.sendErrorResponse(http.StatusBadRequest, -1, err.Error())
			return
		}
		ctx.sendSuccessResponse(res)
	}
}

func JSONRequestBodyWrapper[T any](apiFunc InterfaceWrapFunc[T]) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewHTTPContext(w, r)
		if r.ContentLength <= 0 {
			ctx.sendErrorResponse(http.StatusBadRequest, -1, "missing request body")
			return
		}
		val, err := decodeBody[T](r)
		if err != nil {
			ctx.sendErrorResponse(http.StatusBadRequest, -1, err.Error())
			return
		}
		res, err := apiFunc(val, ctx)
		if err != nil {
			var apiError *errval.ApiError
			if errors.As(err, &apiError) {
				ctx.sendErrorResponse(http.StatusBadRequest, apiError.Code, apiError.Error())
				return
			}
			ctx.sendErrorResponse(http.StatusBadRequest, -1, err.Error())
			return
		}
		ctx.sendSuccessResponse(res)
	}
}

func QueryParamWrapper[T any](apiFunc InterfaceWrapFunc[T]) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewHTTPContext(w, r)
		params, err := decodeQueryParams[T](r)
		if err != nil {
			ctx.sendErrorResponse(http.StatusBadRequest, 400, err.Error())
			return
		}
		res, err := apiFunc(params, ctx)
		if err != nil {
			ctx.sendErrorResponse(http.StatusInternalServerError, -1, err.Error())
			return
		}
		ctx.sendSuccessResponse(res)
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

func decodeQueryParams[T any](r *http.Request) (*T, error) {
	values := new(T)
	if err := decoder.Decode(values, r.URL.Query()); err != nil {
		return nil, err
	}
	if err := validation.ValidateStruct(values); err != nil {
		return nil, err
	}
	return values, nil
}
