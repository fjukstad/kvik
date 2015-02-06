package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/fjukstad/kvik/utils"
	zmq "github.com/pebbe/zmq4"

	"flag"
	"fmt"
)

var workerPort int

func worker(b []byte, filename string) {
	storeScript(b, filename)
	startWorker(filename)
}

func startWorker(filename string) {
	cmd := exec.Command("python", "worker/worker.py", strconv.Itoa(workerPort), filename)
	err := cmd.Start()
	if err != nil {
		log.Panic(err)
	}
}

// storeScript stores the r script from *b* into filename
func storeScript(b []byte, filename string) {

	err := utils.CreateDirectories(filename)
	if err != nil {
		log.Panic(err)
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	_, err = file.Write(b)
	if err != nil {
		log.Panic(err)
	}
	return
}

// Compute is the component of Kvik responsible for starting/stopping workers
// that perform the statistical analyses. It exposes a http rest interface to
// the outside world.
func main() {

	var ip = flag.String("ip", "*", "ip to run on")
	var port = flag.String("port", ":8888", "port to run on")

	flag.Parse()

	responder, _ := zmq.NewSocket(zmq.REP)

	defer responder.Close()

	addr := "tcp://" + *ip + *port

	log.Println(addr)

	err := responder.Bind(addr)

	if err != nil {
		log.Panic(err)
	}

	// ID to identify client
	id := 0

	workerPort = 5000

	for {
		msg, err := responder.Recv(0)

		if err != nil {
			log.Panic(err)
		}

		fmt.Println("Received ", msg)

		cmd := new(utils.Command)
		err = json.Unmarshal([]byte(msg), cmd)
		if err != nil {
			log.Panic(err)
		}

		switch {
		case cmd.Command == "stop":
			break
		case cmd.Command == "startWorker":
			path := "scripts/" + strconv.Itoa(id) + "/script.r"
			worker(cmd.File, path)
		}

		responder.Send(strconv.Itoa(workerPort), 0)

		id += 1
		workerPort += 1
	}

}
