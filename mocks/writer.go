package mocks

type Writer struct {
	WriteFn       func(p []byte) (n int, err error)
	WriteFnCalled bool
}

func NewWriter() *Writer {
	return &Writer{WriteFn: defaultWriteFn}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.WriteFnCalled = true
	return w.WriteFn(p)
}

var defaultWriteFn = func(p []byte) (n int, err error) {
	return
}

func NewTestWriteFn(testFunc func(string)) func(p []byte) (n int, err error) {
	return func(p []byte) (n int, err error) {
		testFunc(string(p))
		return 0, nil
	}
}
