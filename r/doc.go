/*
Package r provides an interface to the R statistical programming language. Users
can interface with R through three methods: Call, Get or Rpc. Call executes a
function call and returns a unique key that can be used in subseqent calls to
Get to retrieve the results. Get can return results in json, csv, png or pdf.
Rpc can be used to bundle Call and Get in a single function call.

As a single R session:

	session, err := r.NewSession()
	...
	key, err := session.Call("stats","rnorm","n=100")
	...
	res, err := session.Get(key, "json")
	...
	res, err = session.Rpc("stats", "rnorm", "n=100", "json")
	...

standalone server that accepts HTTP requests (e.g. from the Client or curl)

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

*/
package r
