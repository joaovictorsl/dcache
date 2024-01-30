package client

import "fmt"

const (
	NOT_ACTIVE_CONN = iota
	CONN_ERROR
	FAILED_TO_CONNECT
	TERMINATED_CLIENT
	KEY_NOT_FOUND
)

type DCacheError struct {
	msg  string
	code uint
}

func dCacheNotActiveConnError(addr string) *DCacheError {
	return &DCacheError{
		msg:  fmt.Sprintf("(%s) connection is not active", addr),
		code: NOT_ACTIVE_CONN,
	}
}

func dCacheConnError(err error) *DCacheError {
	return &DCacheError{
		msg:  err.Error(),
		code: CONN_ERROR,
	}
}

func dCacheFailedToConnectError(addr string, err error) *DCacheError {
	return &DCacheError{
		msg:  fmt.Sprintf("(%s) All attempts failed to connect: %s", addr, err.Error()),
		code: FAILED_TO_CONNECT,
	}
}

func dCacheTerminatedClientError() *DCacheError {
	return &DCacheError{
		msg:  "this client is terminated",
		code: TERMINATED_CLIENT,
	}
}

func dCacheKeyNotFoundError(key string) *DCacheError {
	return &DCacheError{
		msg:  fmt.Sprintf("key (%s) was not found", key),
		code: KEY_NOT_FOUND,
	}
}

func (dcerr *DCacheError) Error() string {
	return dcerr.msg
}

func (dcerr *DCacheError) Code() uint {
	return dcerr.code
}
