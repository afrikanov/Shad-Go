// +build !solution

package otp

import (
	"errors"
	"io"
	"testing/iotest"
)

type cipherReader struct {
	r    io.Reader
	prng io.Reader
}

func (cr cipherReader) Read(p []byte) (n int, e error) {
	size, err := cr.r.Read(p)
	randoms := make([]byte, size)
	n, _ = cr.prng.Read(randoms)
	for i := 0; i < n; i++ {
		p[i] = p[i] ^ randoms[i]
	}
	if n < size {
		err = errors.New("need more bytes")
	}
	return n, err
}

type cipherWriter struct {
	w    io.Writer
	prng io.Reader
}

func (cw cipherWriter) Write(p []byte) (n int, e error) {
	decodeRands := make([]byte, len(p))
	cnt, _ := cw.prng.Read(decodeRands)
	if cnt == 0 {
		return n, iotest.ErrTimeout
	}
	i := 0
	result := make([]byte, cnt)
	for ; i < cnt; i++ {
		result[i] = p[i] ^ decodeRands[i]
	}
	n, err := cw.w.Write(result)
	return n, err
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return cipherReader{r, prng}
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return cipherWriter{w, prng}
}
