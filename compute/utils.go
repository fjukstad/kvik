package compute

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/fjukstad/kvik/utils"
)

// Contacts the compute master to start up a new worker
func StartWorker(addr string, filename string) (string, error) {

	requester, _ := zmq.NewSocket(zmq.REQ)
	//defer requester.Close()

	err := requester.Connect(addr)
	if err != nil {
		return "", err
	}

	f, err := ioutil.ReadFile(filename)
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
