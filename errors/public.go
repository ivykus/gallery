package errors

// Public wraps the original error with a new error that has
// "Public() string" method that will return a message that is
// acceptable to show to the user. This error can also be unwrapped
// using traditional 'errors' package approach.
func Public(err error, msg string) error {
	return publicError{err, msg}
}

type publicError struct {
	err error
	msg string
}

func (p publicError) Error() string {
	return p.err.Error()
}

func (p publicError) Public() string {
	return p.msg
}

func (p publicError) Unwrap() error {
	return p.err
}
