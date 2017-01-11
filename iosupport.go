package filesystem

import (
	"golang.org/x/net/context"
	"io"
)

/*
Implementation of a wrapper for ReadCloser which ignores deadlines and
cancellations to make it compatible with io.ReadCloser.
*/
type ioCompatReadCloser struct {
	readCloser ReadCloser
}

/*
See io.ReadCloser#Read
*/
func (rc *ioCompatReadCloser) Read(p []byte) (int, error) {
	var ctx = context.Background()
	return rc.readCloser.Read(ctx, p)
}

/*
See io.ReadCloser#Close
*/
func (rc *ioCompatReadCloser) Close() error {
	var ctx = context.Background()
	return rc.readCloser.Close(ctx)
}

/*
ReadCloser is a context-aware variant of the good old io.ReadCloser.
*/
type ReadCloser interface {
	Read(context.Context, []byte) (int, error)
	Close(context.Context) error
}

/*
ToIoReadCloser creates a context-ignorant object for providing an
io.ReadCloser compatible API.
*/
func ToIoReadCloser(rc ReadCloser) io.ReadCloser {
	return &ioCompatReadCloser{readCloser: rc}
}

/*
Implementation of a wrapper for WriteCloser which ignores deadlines and
cancellations to make it compatible with io.WriteCloser.
*/
type ioCompatWriteCloser struct {
	writeCloser WriteCloser
}

/*
See io.WriteCloser#Write
*/
func (wc *ioCompatWriteCloser) Write(p []byte) (int, error) {
	var ctx = context.Background()
	return wc.writeCloser.Write(ctx, p)
}

/*
See io.WriteCloser#Close
*/
func (wc *ioCompatWriteCloser) Close() error {
	var ctx = context.Background()
	return wc.writeCloser.Close(ctx)
}

/*
WriteCloser is a context-aware variant of the good old io.WriteCloser.
*/
type WriteCloser interface {
	Write(context.Context, []byte) (n int, err error)
	Close(context.Context) error
}

/*
ToIoWriteCloser creates a context-ignorant object for providing a
io.WriteCloser compatible API.
*/
func ToIoWriteCloser(wc WriteCloser) io.WriteCloser {
	return &ioCompatWriteCloser{writeCloser: wc}
}
