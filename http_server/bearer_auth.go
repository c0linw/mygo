package http_server

import (
	"errors"
	"mygo/errval"
	"net/http"
	"strings"
)

func ProcessBearer(tokenFunc func(token string) error) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := NewHTTPContext(w, r)
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				ctx.sendErrorResponse(http.StatusUnauthorized, -1, "no token")
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				ctx.sendErrorResponse(http.StatusUnauthorized, -1, "invalid token format")
				return
			}
			idToken := parts[1]
			err := tokenFunc(idToken)
			if err != nil {
				var apiError *errval.ApiError
				if errors.As(err, &apiError) {
					ctx.sendErrorResponse(http.StatusUnauthorized, apiError.Code, apiError.Error())
				}
				ctx.sendErrorResponse(http.StatusUnauthorized, -1, err.Error())
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
