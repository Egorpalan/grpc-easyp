package exception

import (
	"errors"

	pb "github.com/Egorpalan/grpc-easyp/pkg/api/notes/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	ErrNoteNotFound = errors.New("note not found")
	ErrInvalidInput = errors.New("invalid input")
)

var errorMap = map[error]codes.Code{
	ErrNoteNotFound: codes.NotFound,
	ErrInvalidInput: codes.InvalidArgument,
}

var errorCodeMap = map[error]pb.ErrorCode{
	ErrNoteNotFound: pb.ErrorCode_ERROR_CODE_NOTE_NOT_FOUND,
	ErrInvalidInput: pb.ErrorCode_ERROR_CODE_INVALID_INPUT,
}

func Get(err error) codes.Code {
	code, ok := errorMap[err]
	if ok {
		return code
	}

	for e, errorCode := range errorMap {
		if errors.Is(err, e) {
			return errorCode
		}
	}
	return codes.Internal
}

func GetErrorCode(err error) pb.ErrorCode {
	code, ok := errorCodeMap[err]
	if ok {
		return code
	}

	for e, errorCode := range errorCodeMap {
		if errors.Is(err, e) {
			return errorCode
		}
	}
	return pb.ErrorCode_ERROR_CODE_UNSPECIFIED
}

func WrapError(err error) error {
	return status.Error(Get(err), err.Error())
}

func WrapErrorWithDetails(err error, reason string, internalErrorCode string) error {
	code := Get(err)
	st := status.New(code, err.Error())

	customError := &pb.CustomError{
		Code:              GetErrorCode(err),
		Reason:            reason,
		InternalErrorCode: internalErrorCode,
	}

	anyDetails, marshalErr := anypb.New(customError)
	if marshalErr != nil {
		return st.Err()
	}

	st, _ = st.WithDetails(anyDetails)
	return st.Err()
}
