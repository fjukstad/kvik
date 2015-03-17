package dataset

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/fjukstad/kvik/utils"
	zmq "github.com/pebbe/zmq4"
)

type Dataset struct {
	worker *zmq.Socket
}

func (d *Dataset) Call(args ...interface{}) (interface{}, error) {
	req := utils.WorkerRequest{"command", args}

	enc, err := json.Marshal(req)
	if err != nil {
		log.Println("json error", err)
		return "", err
	}

	msg := strings.Replace(string(enc), "\\\"", "'", -1)

	fmt.Println("Sending", msg, "to", d)
	d.worker.Send(msg, 0)
	resp := new(utils.WorkerResponse)

	reply, err := d.worker.Recv(0)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(reply), resp)
	if err != nil {
		log.Println("json error", err)
		return "", err
	}

	if strings.Contains(reply, "Could not evaluate") {
		return "", errors.New(reply)
	}

	if resp.Status != 0 {
		fmt.Println("REPLY:", reply)
		return "", errors.New(reply)
	}

	return resp.Response, nil

}

func ping(d *Dataset) error {
	req := utils.WorkerRequest{"ping", ""}
	enc, err := json.Marshal(req)
	if err != nil {
		fmt.Println("json error", err)
		return err
	}

	d.worker.Send(string(enc), 0)
	resp := new(utils.WorkerResponse)

	reply, err := d.worker.Recv(0)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(reply), resp)
	if err != nil {
		log.Println("json error", err)
		return err
	}

	if !strings.Contains(reply, "pong") {
		return errors.New("Did not receive correct response from worker")
	}
	return nil
}

func ConnectToRunningWorker(addr string) (d *Dataset, err error) {
	requester, _ := zmq.NewSocket(zmq.REQ)
	err = requester.Connect(addr)
	if err != nil {
		fmt.Println("Could not connect to worker", err)
		return nil, err
	}
	d = new(Dataset)
	d.worker = requester

	err = ping(d)

	if err != nil {
		fmt.Println("Could not connect to worker", err)
		return nil, err
	}

	fmt.Println("Connected to worker.")
	return d, nil
}

func RequestNewWorker(ip, port, filename string) (d *Dataset,
	workerAddr string, err error) {

	addr := "tcp://" + ip + port

	// Get a new worker that can do computation for us
	workerPort, err := utils.StartWorker(addr, filename)
	if err != nil {
		fmt.Println("Could not start worker..", workerPort, err)
		return nil, "", err
	}
	workerAddr = "tcp://" + ip + ":" + workerPort

	// Connect to the worker so that we're good to go
	requester, _ := zmq.NewSocket(zmq.REQ)
	err = requester.Connect(workerAddr)
	if err != nil {
		fmt.Println("Could not connect to worker")
		return nil, "", err
	}
	d = new(Dataset)
	d.worker = requester

	err = ping(d)

	if err != nil {
		fmt.Println("Could not connect to worker", err)
		return nil, "", err
	}

	fmt.Println("Worker stared at", workerAddr)
	return d, workerAddr, nil
}

// Contacts the compute master to start up a new worker
func StartWorker(addr string) (string, error) {

	requester, _ := zmq.NewSocket(zmq.REQ)
	//defer requester.Close()

	err := requester.Connect(addr)
	if err != nil {
		return "", err
	}

	f, err := ioutil.ReadFile("scripts/script.r")
	if err != nil {
		return "", err
	}

	// Send start worker message to the master
	c := utils.Command{0, "startWorker", f}
	msg, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	requester.Send(string(msg), 0)

	// Get reply which contains the port where the new worker runs
	reply, err := requester.Recv(0)
	if err != nil {
		return "", err
	}

	worker := new(utils.ComputeResponse)

	err = json.Unmarshal([]byte(reply), worker)
	if err != nil {
		return "", err
	}

	if worker.Status != 1 {
		err = errors.New("Worker not started")
		return "", err
	}

	return worker.Message, nil
}
