package kecpfakews

type FakeWriter struct {
	reliable bool
}

func (fw *FakeWriter) Write(p []byte) (n int, err error) {
	probability := MathRandGen()
	if !fw.reliable && probability < 2 {
		return 0, FakeError
	} else {
		return len(p), nil
	}
}

func (fw *FakeWriter) Close() error {
	probability := MathRandGen()
	if !fw.reliable && probability < 2 {
		return FakeError
	} else {
		return nil
	}
}
