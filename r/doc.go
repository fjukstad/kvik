/*
Package r provides an interface to the R statistical programming language. It
can either be run as a single R session:

	session, err := r.NewSession()
	...
	key, err := session.Call("stats","rnorm","n=100")
	...
	res, err := session.Get(key, "json")
	...
	res, err = session.Rpc("stats", "rnorm", "n=100", "json")
	...

standalone server that accepts HTTP requests (e.g. from the Client):

	s, err := r.Server(4, "/tmp/kvik")
	...
	s.Start(":8181")

or a client if there is a server running somewhere:

	client := r.Client{"http://localhost.com:8181", "username","password"}
	key, err := client.Call("stats","rnorm","n=100")
	...
	res, err := client.Get(key, "json")
	...
	res, err := client.Rpc("stats","rnorm","n=100","json")
	...

The server starts up a number of R sessions that it communicates with to execute
R functions and retrieve results. The server exposes a HTTP interface with three
methods Call, Get and Rpc. Call executes an R function and returns a temporary key
that can be used in subsequent Get requests. Get requests can return json, csv,
png or pdf files. Rpc combines a call to Call following a Get to get the
results.

*/
package r
