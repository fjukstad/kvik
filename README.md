# rpcman

RPC between a Go client and Python server using ØMQ.

Needed somewhere to store this example. 

# Dependencies 
[ØMQ](http://zeromq.org/) and the bindings for
[go](http://zeromq.org/bindings:go) and
[python](http://zeromq.org/bindings:python). 

# How to 
Install the go package
    
    go get github.com/fjukstad/rpcman


Test it by cloning down the repo and running

    python server.py

and 

    go run test/test.go 

# TODO
Make a python lib out of it
