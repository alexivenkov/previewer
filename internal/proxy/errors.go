package proxy

import "errors"

var (
	ErrBadResponse       = errors.New("remote server unavailable")
	ErrMaxSizeExceed     = errors.New("image max size exceed")
	ErrUnsupportedFormat = errors.New("only jpeg or png formats available")
)

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}
