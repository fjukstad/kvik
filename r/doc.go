/*
Package r provides an interface to the R statistical programming language. It
can either be run as a standalone server that accepts HTTP requests:

	s, err := r.Server(4, "/tmp/kvik")
	...
	s.Start(":8181")

or a single R session

	session, err := r.NewSession()
	...
	key, err := session.Call("stats","rnorm","n=100")
	...
	res, err := session.Get(key, "json")
	...

or a client if there is a server running elsewhere:

	client := r.Client{"http://example.com:8181", "username","password"}
	key, err := client.Call("stats","rnorm","n=100")
	...
	res, err := client.Get(key, "json")
	...

The server starts up a number of R sessions that it communicates with to execute
R functions and retrieve results. The server exposes a HTTP interface with two
methods Call and Get. Call executes an R function and returns a temporary key
that can be used in subsequent Get requests. Get requests can return json, csv,
png or pdf files.

*/
package r
