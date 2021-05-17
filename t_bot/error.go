package t_bot

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrorHomeNotFound = errors.New("home location is not found")
	ErrorAlreadyExist = errors.New("home location is already exist")
	ErrorInternal= errors.New("some error, please try again")
)

func fromGRPCErr(err error) error {
	log.Error(err)
	st, _ := status.FromError(err)
	switch st.Code() {
	case codes.NotFound:
		return ErrorHomeNotFound
	default:
		return ErrorHomeNotFound
	}
}
