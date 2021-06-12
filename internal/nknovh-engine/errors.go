package nknovh_engine

type error interface {
	Error() string
}

func createErr(text string) error {
	return &errorString{text}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
