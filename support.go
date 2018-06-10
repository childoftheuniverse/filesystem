package filesystem

import (
	"context"
	"errors"
	"net/url"
)

/*
EUNSUPP is the error returned when an operation is called on a file system
which does not have the requested functionality.
*/
var EUNSUPP = errors.New("File system does not support this operation")

/*
ENOFS is the error which should only ever be returned directly from this
framework indicating that no file system is registered under that name.
*/
var ENOFS = errors.New("No file system loaded for URL type")

/*
FileWatchFunc describes the format of a function which can be used to watch
for changes in a file in a supported file system. The function will be called
when a watched file is modified. The first parameter will be the path of the
file being changed, and the second parameter is a ReadCloser which can be
used to access the modified contents.

Implementations must allow for the ReadCloser to be discarded without ever
calling Read().
*/
type FileWatchFunc func(*url.URL, ReadCloser)

/*
CancelWatchFunc describes a function to call to stop watching a file.
*/
type CancelWatchFunc func() error

/*
FileSystem provides the expected minimal interface required from a file system
implementation.

File systems should support at least a subset of these functions; any
functions which are unsupported must return an EUNSUPP error.
Deadlines and cancellations should be supported by either the underlying
file system or implemented as reasonably as possible in the adapter.
*/
type FileSystem interface {
	// Open the specified file for reading. Context should be used to control
	// the opening, not the reading part of the operation.
	OpenReader(context.Context, *url.URL) (ReadCloser, error)

	// Open the specified file for writing. The affected file should be
	// overwritten. Context should be used to control the opening, not the
	// writing part of the operation.
	OpenWriter(context.Context, *url.URL) (WriteCloser, error)

	// Open the specified file for appending. The affected file should be
	// created if it did not exist, but no contents should be overwritten.
	// Context should be used to control the opening, not the writing part
	// of the operation.
	OpenAppender(context.Context, *url.URL) (WriteCloser, error)

	// Retrieve a list of entries described by the URL, if the file system
	// has such a notion. Should return the relative file names, without
	// ., .. or whatever their equivalent is.
	ListEntries(context.Context, *url.URL) ([]string, error)

	// Watch for changes in a given file and call the FileWatchFunc on every
	// change to the watched file. Watching anything other than files is left
	// as an implementation detail.
	WatchFile(context.Context, *url.URL, FileWatchFunc) (CancelWatchFunc, chan error, error)

	// Delete the specified file. Failures may or may not leave the file
	// existing.
	Remove(context.Context, *url.URL) error
}

/*
All file system implementation adapters will be registered in this map.
*/
var registeredFileSystems = make(map[string]FileSystem)

/*
AddImplementation is used on initialization of individual file system modules
to sign file systems up for receiving calls through the API. Any calls to
URLs with the given scheme will be patched through to the file system
implementation.

Subsequent invocations of AddImplementation will cause the association to be
overwritten.

This function may be called from init() for easy file systems, or may require
a more involved setup procedure for file systems talking to a server node
and/or requiring authentication.
*/
func AddImplementation(scheme string, fs FileSystem) {
	registeredFileSystems[scheme] = fs
}

/*
GetImplementation fetches a pointer to the entire implementation of the file
system which would be used to handle the URL. If no file system can handle
the URL, this returns nil.

Usually you will want to use one of the more specific functions.
*/
func GetImplementation(fileurl *url.URL) FileSystem {
	var found bool

	if _, found = registeredFileSystems[fileurl.Scheme]; found {
		return registeredFileSystems[fileurl.Scheme]
	}

	return nil
}

/*
HasImplementation determines whether a handler is registered for the given
schema. Returns true if an implementation was registered for the specified
scheme.
*/
func HasImplementation(scheme string) bool {
	var found bool

	_, found = registeredFileSystems[scheme]

	return found
}

/*
OpenReader opens the referenced file and returns a ReadCloser object which
can be used to access the files contents.
*/
func OpenReader(ctx context.Context, fileurl *url.URL) (ReadCloser, error) {
	var fs = GetImplementation(fileurl)

	if fs == nil {
		return nil, ENOFS
	}

	return fs.OpenReader(ctx, fileurl)
}

/*
OpenWriter opens the referenced file and returns a WriteCloser object which
can be used to put data into the file. Any previous file contents will be
overwritten.

Implementations may require Close() to be invoked before any changes are made
whatsoever.
*/
func OpenWriter(ctx context.Context, fileurl *url.URL) (WriteCloser, error) {
	var fs = GetImplementation(fileurl)

	if fs == nil {
		return nil, ENOFS
	}

	return fs.OpenWriter(ctx, fileurl)
}

/*
OpenAppender opens the referenced file and returns a WriteCloser object which
can be used to append data to the file. If the file does not exist, it will
be created.

Implementations may require Close() to be invoked before any changes are made
whatsoever.
*/
func OpenAppender(ctx context.Context, fileurl *url.URL) (WriteCloser, error) {
	var fs = GetImplementation(fileurl)

	if fs == nil {
		return nil, ENOFS
	}

	return fs.OpenAppender(ctx, fileurl)
}

/*
ListEntries retrieves a list of relative object names beneath the specified
URL. Objects may be something resembling to files or directories and will not
contain special entries such as the local and parent directory.
*/
func ListEntries(ctx context.Context, dirurl *url.URL) ([]string, error) {
	var fs = GetImplementation(dirurl)

	if fs == nil {
		return nil, ENOFS
	}

	return fs.ListEntries(ctx, dirurl)
}

/*
WatchFile waits for modifications of the file at the specified URL and invokes
the watcher with any modified files. Some implementations may allow
watching directories.
*/
func WatchFile(ctx context.Context, fileurl *url.URL, watcher FileWatchFunc) (
	CancelWatchFunc, chan error, error) {
	var fs = GetImplementation(fileurl)

	if fs == nil {
		return nil, nil, ENOFS
	}

	return fs.WatchFile(ctx, fileurl, watcher)
}

/*
Remove deletes the referenced object from the underlying file system.
Removal is guaranteed to succeed if no error returns, otherwise it may or may
not have succeeded.
*/
func Remove(ctx context.Context, fileurl *url.URL) error {
	var fs = GetImplementation(fileurl)

	if fs == nil {
		return ENOFS
	}

	return fs.Remove(ctx, fileurl)
}
