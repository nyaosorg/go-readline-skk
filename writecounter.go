package skk

type writeCounter struct {
	n   int64
	err error
}

func (w *writeCounter) Try(n int, err error) bool {
	w.n += int64(n)
	w.err = err
	return err != nil
}

func (w *writeCounter) Try64(n int64, err error) bool {
	w.n += n
	w.err = err
	return err != nil
}

func (w *writeCounter) Result() (int64, error) {
	return w.n, w.err
}
