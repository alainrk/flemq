package common

type OffsetNotFoundError struct {
	Err error
}

func (e OffsetNotFoundError) Error() string {
	return e.Err.Error()
}
