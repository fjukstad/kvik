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

// workerPort holds the next available port to start up a worker.
var workerPort int

func worker(b []byte, filename string) error {
	err := storeScript(b, filename)
	if err != nil {
		return err
	}
	err = startWorker(filename)
	if err != nil {
		return err
	}
	return nil
}

func startWorker(filename string) error {
	cmd := exec.Command("python", "worker/worker.py", strconv.Itoa(workerPort),
		filename)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

// storeScript stores the r script from *b* into filename
func storeScript(b []byte, filename string) error {
	err := utils.CreateDirectories(filename)
	if err != nil {
		log.Panic(err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func killAllWorkers() error {
	cmd := exec.Command("killall", "Python")
	return cmd.Start()
}

func errorResponse(msg string) utils.ComputeResponse {
	return utils.ComputeResponse{-1, msg}
}
func successResponse(msg string) utils.ComputeResponse {
	return utils.ComputeResponse{1, msg}
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
			log.Println(err)
			continue
		}
		fmt.Println("Received ", msg)

		cmd := new(utils.Command)
		err = json.Unmarshal([]byte(msg), cmd)
		if err != nil {
			log.Panic(err)
		}

		var Response utils.ComputeResponse

		switch {
		case cmd.Command == "stop":
			break

		case cmd.Command == "startWorker":
			path := "scripts/" + strconv.Itoa(id) + "/script.r"
			err = worker(cmd.File, path)
			if err != nil {
				log.Println(err)
				Response = errorResponse("Could not start worker")
			} else {
				Response = successResponse(strconv.Itoa(workerPort))
			}

		case cmd.Command == "killall":
			err = killAllWorkers()
			if err != nil {
				Response = errorResponse("Could not kill all workers")
			} else {
				Response = successResponse("Killed off all workers")
			}
		}

		resp, _ := json.Marshal(Response)
		responder.Send(string(resp), 0)
		id += 1
		workerPort += 1
	}
}
