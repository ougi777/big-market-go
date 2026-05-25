package types

type Response[T any] struct {
	Code string `json:"code"`
	Info string `json:"info"`
	Data T      `json:"data,omitempty"`
}

func Success[T any](data T) Response[T] {
	return Response[T]{
		Code: ResponseCodeSuccess.Code,
		Info: ResponseCodeSuccess.Info,
		Data: data,
	}
}

func Failure[T any](code ResponseCode, data T) Response[T] {
	return Response[T]{
		Code: code.Code,
		Info: code.Info,
		Data: data,
	}
}
