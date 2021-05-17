package t_bot

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrorHomeNotFound = errors.New("home location is not found")
	ErrorAlreadyExist = errors.New("home location is already exist")
	ErrorInternal= errors.New("some error, please try again")
)

func fromGRPCErr(err error) error {
	st, _ := status.FromError(err)
	switch st.Code() {
	case codes.NotFound:
		return ErrorHomeNotFound
	default:
		return ErrorHomeNotFound
	}
}
