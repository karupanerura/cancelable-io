package cio

import (
	"context"
	"io"
)

type CalcelableReader struct {
	ctx context.Context
	r   io.Reader
}

type readResult struct {
	n   int
	err error
}

func NewCancelableReader(ctx context.Context, r io.Reader) *CalcelableReader {
	return &CalcelableReader{ctx, r}
}

func (c *CalcelableReader) Read(b []byte) (int, error) {
	res := make(chan readResult)

	go func() {
		n, err := c.r.Read(b)
		res <- readResult{n, err}
		close(res)
	}()

	select {
	case result := <-res:
		return result.n, result.err
	case <-c.ctx.Done():
		return 0, context.DeadlineExceeded
	}
}
