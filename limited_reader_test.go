package filesystem

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"testing"
)

var ErrExpected = errors.New("Expect this error")

type MockReadCloser struct {
	Fail bool
}

func (r *MockReadCloser) Read(ctx context.Context, p []byte) (int, error) {
	if r.Fail {
		return 0, ErrExpected
	}
	return rand.Read(p)
}

func (r *MockReadCloser) Close(ctx context.Context) error {
	if r.Fail {
		return ErrExpected
	}
	return nil
}

func TestLimitedReadCloserReadTooMuch(t *testing.T) {
	buf := make([]byte, 100)
	mockReadCloser := &MockReadCloser{}

	l := LimitedReadCloser{R: mockReadCloser, N: 25}

	n, err := l.Read(context.Background(), buf)

	if err != nil {
		t.Errorf("Error reported from read call: %s", err.Error())
	}

	if n != 25 {
		t.Errorf("Read reported wrong length %d", n)
	}

	n, err = l.Read(context.Background(), buf)
	if err != io.EOF {
		t.Error("Second read did not return EOF")
	}
	if n != 0 {
		t.Errorf("Second read length was %d, expected 0", n)
	}
}

func TestLimitedReadCloserMultipleReads(t *testing.T) {
	buf := make([]byte, 100)
	mockReadCloser := &MockReadCloser{}

	l := LimitedReadCloser{R: mockReadCloser, N: 125}

	n, err := l.Read(context.Background(), buf)

	if err != nil {
		t.Errorf("Error reported from first read call: %s", err.Error())
	}

	if n != 100 {
		t.Errorf("Read reported wrong length %d (expected 100)", n)
	}

	n, err = l.Read(context.Background(), buf)

	if err != nil {
		t.Errorf("Error reported from second read call: %s", err.Error())
	}

	if n != 25 {
		t.Errorf("Read reported wrong length %d (expected 25)", n)
	}

	n, err = l.Read(context.Background(), buf)
	if err != io.EOF {
		t.Errorf("Second read did not return EOF (%v)", err)
	}
	if n != 0 {
		t.Errorf("Second read length was %d, expected 0", n)
	}
}

func TestLimitedReadCloserForwardsReadError(t *testing.T) {
	buf := make([]byte, 100)
	mockReadCloser := &MockReadCloser{Fail: true}

	l := LimitedReadCloser{R: mockReadCloser, N: 125}

	_, err := l.Read(context.Background(), buf)

	if err != ErrExpected {
		t.Errorf("Unexpected error from Read: %v", err)
	}
}

func TestLimitedReadCloserForwardsCloseError(t *testing.T) {
	mockReadCloser := &MockReadCloser{Fail: true}

	l := LimitedReadCloser{R: mockReadCloser, N: 1}

	err := l.Close(context.Background())

	if err != ErrExpected {
		t.Errorf("Unexpected error from Close: %v", err)
	}
}
