package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

func StartWorkerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	log.Print(body)
}

// Compute is the component of Kvik responsible for starting/stopping workers
// that perform the statistical analyses. It exposes a http rest interface to
// the outside world.
func main() {

	var ip = flag.String("ip", "", "ip to run on")
	var port = flag.String("port", ":8888", "port to run on")

	flag.Parse()

	log.Print("Compute master ", *ip, *port)

	http.HandleFunc("/startworker", StartWorkerHandler)

	err = http.ListenAndServe(*port, nil)
	if err != nil {
		log.Print(err)
	}

}
