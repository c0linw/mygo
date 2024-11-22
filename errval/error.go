package errval

type ApiError struct {
	Code    int
	Message string
}

func NewError(code int, msg string) *ApiError {
	return &ApiError{Code: code, Message: msg}
}

func (e *ApiError) Error() string {
	return e.Message
}
