/*
filesystem provides an abstraction API for hooking up various file system
implementations and use them more or less transparently through URLs and a
common API.

In this API, the scheme part of the URL (see net/url for a description of what
makes up an URL) is used to find the correct file system implementation. The
path part should be used to point to the correct file, while the query part
may be used to supply additional metadata for the operation.

Unlike the regular go io implementation, an effort is made by this API to
provide context functionality (deadlines, cancellation) in all relevant
functions. Implementations of file system adapters should make a point of
providing these semantics in a reasonable way; refer to the information in
README.md for more details.
*/
package filesystem
