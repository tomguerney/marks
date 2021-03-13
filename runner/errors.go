package runner

type runnerError struct {
	msg string
}

func (e runnerError) Error() string {
	return e.msg
}
