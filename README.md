# Generalized file system access for golang projects

filesystem mostly offers an abstract API which can be used to adapt actual
file implementations to. filesystem provides a small common denominator of
access routines for file systems which will be sufficient for most basic
applications.

Actual implementations of file systems will be provided in separate projects
and can be added in a plugin-like model.

## Using the abstraction API

Actual implementations of file system adapters must be loaded before any of
the functions can be used. Simple file systems can be registered by just
loading the module (as an anonymous import), more complex file systems may
require invoking initialization functions for authenticating or pointing to
some kind of frontend.

The abstraction API is mostly contained in the directly exported functions
in the support.go module, such as OpenReader, OpenWriter, ListEntries, etc.
They should be relatively well documented and self-explanatory.

The read/write API mostly uses internal ReadCloser/WriteCloser abstractions
instead of the regular io counterparts to preserve operation contexts, which
are (supposed to be) a major part of the idea of this API. If you need to use
regular io.ReadCloser/io.WriteCloser operations on these (ignoring deadlines
etc.), use ToIoReadCloser()/ToIoWriteCloser() to get the corresponding
representations.

## Implementing file system adapters

If you want to implement adapters for existing file system APIs, simply
implement all methods of the FileSystem object in support.go in your file
system object and call AddImplementation() on the resulting object. It is
entirely reasonable to only call AddImplementation() in some kind of
initialization routine which handles authentication or some other kind
of configuration; if your file system does not need such a thing, just call
it from init().

Please keep all the intended semantics of the context APIs intact as much as
possible. This would mean, for example, that a Read() which exceeds its
deadline should probably attempt to restore the previous position pointer in
the file to the place it was at before. A Write() should probably truncate the
file to its previous length upon exceeding deadlines or cancellation.

Ultimately, expectations of the API users should be kept in mind and used as
a guideline to make this API most useful.
