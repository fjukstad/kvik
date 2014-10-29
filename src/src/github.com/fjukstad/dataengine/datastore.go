package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/fjukstad/rpcman"
)

var rpc *rpcman.RPCMan

func main() {

	var path = flag.String("path", "/Users/bjorn/stallo/data", "path where data files are stored")
	var ip = flag.String("ip", "localhost", "ip to run on")
	var port = flag.String("port", ":8888", "port to run on")

	flag.Parse()

	log.Print("Starting datastore at ", *ip, *port)

	http.HandleFunc("/com", CommandHandler)

	http.ListenAndServe(*port, nil)

}
