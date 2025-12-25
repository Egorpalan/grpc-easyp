package exception

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNoteNotFound = errors.New("note not found")
	ErrInvalidInput = errors.New("invalid input")
)

var errorMap = map[error]codes.Code{
	ErrNoteNotFound: codes.NotFound,
	ErrInvalidInput: codes.InvalidArgument,
}

func Get(err error) codes.Code {
	code, ok := errorMap[err]
	if ok {
		return code
	}

	for e, code := range errorMap {
		if errors.Is(err, e) {
			return code
		}
	}
	return codes.Internal
}

func WrapError(err error) error {
	return status.Error(Get(err), err.Error())
}
