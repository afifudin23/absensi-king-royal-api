package common

type ErrorSchema struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details any       `json:"details"`
}

type Response[T any] struct {
	Success bool         `json:"success"`
	Data    *T           `json:"data"`
	Error   *ErrorSchema `json:"error"`
}

func SuccessResponse[T any](data T) Response[T] {
	return Response[T]{
		Success: true,
		Data:    &data,
		Error:   nil,
	}
}

func ErrorResponse[T any](errSchema ErrorSchema) Response[T] {
	return Response[T]{
		Success: false,
		Data:    nil,
		Error:   &errSchema,
	}
}
