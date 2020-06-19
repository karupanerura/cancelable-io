package cio

import (
	"context"
	"io"
)

type CalcelableWriter struct {
	ctx context.Context
	w   io.Writer
}

type writeResult struct {
	n   int
	err error
}

func NewCancelableWriter(ctx context.Context, w io.Writer) *CalcelableWriter {
	return &CalcelableWriter{ctx, w}
}

func (c *CalcelableWriter) Write(b []byte) (int, error) {
	res := make(chan writeResult)

	go func() {
		n, err := c.w.Write(b)
		res <- writeResult{n, err}
		close(res)
	}()

	select {
	case result := <-res:
		return result.n, result.err
	case <-c.ctx.Done():
		return 0, context.DeadlineExceeded
	}
}
