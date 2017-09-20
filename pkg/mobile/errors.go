package mobile

type NotFoundError struct {
	Message string
}

func (nfe *NotFoundError) Error() string {
	return nfe.Message
}

func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

type StatusError struct {
	Message string
	Code    int
}

func (se *StatusError) Error() string {
	return se.Message
}

func (se *StatusError) StatusCode() int {
	if se.Code == 0 {
		return se.DefaultStatusCode()
	}
	return se.Code
}

func (se *StatusError) DefaultStatusCode() int {
	return 500
}
