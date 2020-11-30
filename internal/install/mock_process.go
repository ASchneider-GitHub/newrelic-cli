package install

type mockProcess struct {
	name string
	pid  int32
}

func (p mockProcess) Name() (string, error) {
	return p.name, nil
}

func (p mockProcess) PID() int32 {
	return p.pid
}
