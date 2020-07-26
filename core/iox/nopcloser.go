package iox

import "io"

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error {
	return nil
}

func NopCloser(w io.Writer) io.WriteCloser {
	return nopCloser{w}
}
