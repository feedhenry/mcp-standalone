package data

type notFoundError struct {
	Message string
	Code    int
}

func (nf *notFoundError) Error() string {
	return nf.Message
}
func (nf *notFoundError) ErrorCode() int {
	return nf.Code
}

func newNotFoundError(message string) *notFoundError {
	return &notFoundError{
		Message: message,
		Code:    404,
	}
}

func IsNotFoundErr(err error) bool {
	_, ok := err.(*notFoundError)
	return ok
}
