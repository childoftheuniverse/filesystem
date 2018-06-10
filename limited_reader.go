package filesystem

import (
	"context"
	"io"
)

/*
LimitedReadCloser is the equivalent to io.LimitedReader for the filesystem API.
*/
type LimitedReadCloser struct {
	/* Underlying reader object */
	R ReadCloser

	/* Maximum length to read from the reader before feigning EOF. */
	N int64
}

/*
Read reads up to cap(p) bytes from the underlying file system object.
Make sure the total number of bytes read from the wrapped reader does
not exceed N.
*/
func (l *LimitedReadCloser) Read(ctx context.Context, p []byte) (int, error) {
	var n int
	var err error

	if l.N <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > l.N {
		p = p[0:l.N]
	}
	n, err = l.R.Read(ctx, p)
	l.N -= int64(n)
	return n, err
}

/*
Close calls the close method of the underlying implementation.
*/
func (l *LimitedReadCloser) Close(ctx context.Context) error {
	return l.R.Close(ctx)
}
