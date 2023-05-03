package helper

import (
	"io"
)

// passThru is inspired by https://stackoverflow.com/a/25645804
type passThru struct {
	io.Reader
	current    int64
	length     int64
	percentage chan int
}

func NewProgressReader(r io.ReadCloser, length int64, percentage chan int) io.Reader {
	return &passThru{
		Reader:     r,
		length:     length,
		percentage: percentage,
	}
}

func (pt *passThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	if n > 0 {
		pt.current += int64(n)
		pt.percentage <- int(float64(pt.current) / float64(pt.length) * float64(100))
	}
	return n, err
}
