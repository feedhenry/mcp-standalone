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
