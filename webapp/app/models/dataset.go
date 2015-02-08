package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/fjukstad/kvik/utils"
	zmq "github.com/pebbe/zmq4"
	"github.com/revel/revel"
)

type Dataset struct {
	worker *zmq.Socket
}

// Retrieve Gene expression for the given genes
func (d *Dataset) Exprs(genes []string) []string {
	geneVector := utils.ListToVector(genes)
	command := "exprs(" + geneVector + ")"

	resp, err := d.Call(command)
	if err != nil {
		log.Panic(err)
	}

	response := prepareResponse(resp)
	return response
}

func prepareResponse(resp interface{}) []string {

	var response string
	t := reflect.TypeOf(resp).String()

	if t == "float64" {
		res := resp.(float64)
		response = strconv.FormatFloat(res, 'f', 9, 64)
	}
	if t == "string" {
		response = resp.(string)
	}

	response = strings.Trim(response, "[1] ") // remove r thing
	response = strings.Trim(response, "\n")   // unwanted newlines
	response = strings.TrimLeft(response, " ")
	results := strings.Split(response, " ")

	var res []string
	for _, r := range results {
		if r != "" {
			res = append(res, r)
		}
	}

	return res
}

// Retrieve fold change and associated p-value for the genes
func (d *Dataset) Fc(genes []string) []string {
	geneVector := utils.ListToVector(genes)
	command := "fc(" + geneVector + ")"

	resp, err := d.Call(command)
	if err != nil {
		log.Panic(err)
	}

	response := prepareResponse(resp)
	return response
}

func (d *Dataset) Call(args ...interface{}) (interface{}, error) {
	req := utils.WorkerRequest{"command", args}

	enc, err := json.Marshal(req)
	if err != nil {
		log.Println("json error", err)
		return "", err
	}

	fmt.Println("Sending", string(enc))

	d.worker.Send(string(enc), 0)
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

	if resp.Status != 0 {
		fmt.Println("REPLY:", reply)
		return "", errors.New(reply)
	}

	return resp.Response, nil

}

func InitDataset() (*Dataset, error) {
	ip, _ := revel.Config.String("compute.ip")
	port, _ := revel.Config.String("compute.port")
	addr := ip + port

	// Get a new worker that can do computation for us
	workerPort, err := StartWorker(addr)
	if err != nil {
		return nil, err
	}
	workerAddr := ip + ":" + workerPort

	// Connect to the worker so that we're good to go
	requester, _ := zmq.NewSocket(zmq.REQ)
	//defer requester.Close()
	err = requester.Connect(workerAddr)
	if err != nil {
		return nil, err
	}
	d := new(Dataset)
	d.worker = requester

	return d, nil
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
