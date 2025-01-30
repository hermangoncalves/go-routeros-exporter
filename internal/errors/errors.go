package errors

import "errors"

var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrDeviceUnreachable  = errors.New("device unreachable")
    ErrNetworkTimeout     = errors.New("network timeout")
)