package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

	zmq "github.com/pebbe/zmq4"
)

// Command to compute master
type Command struct {
	Type    int
	Command string
	File    []byte
}

// Response from compute master
type ComputeResponse struct {
	Status  int
	Message string
}

// Requests to compute workers
type WorkerRequest struct {
	Method string
	Args   interface{}
}

// Response from compute workers
type WorkerResponse struct {
	Response interface{}
	Status   int
}

// Used to send results from statistical analyses back to the client
type ClientCompResponse struct {
	Output map[string]interface{}
}

type SearchResponse struct {
	Terms []string
}

func CreateDirectories(filename string) error {

	filepath := path.Dir(filename)
	directories := strings.Split(filepath, "/")

	p := ""
	for i := range directories {
		p += directories[i] + "/"
		err := os.Mkdir(p, 0755)

		if err != nil {
			pe, _ := err.(*os.PathError)

			// if folder exists, continue to the next one
			if !strings.Contains(pe.Error(), "file exists") {
				log.Println("Mkdir failed miserably: ", err)
				return err
			}
		}
	}

	return nil
}

func ListToVector(list []string) string {
	str := "c("
	for i, b := range list {
		if i < 1 {
			str = str + "\"" + b + "\""
		} else {
			str = str + "," + "\"" + b + "\""
		}
	}
	str = str + ")"
	return str
}

// Contacts the compute master to start up a new worker
func StartWorker(addr string, filename string) (string, error) {

	requester, _ := zmq.NewSocket(zmq.REQ)
	//defer requester.Close()

	err := requester.Connect(addr)
	if err != nil {
		fmt.Println("Coud not connect to compute at ", addr)
		return "", err
	}

	f, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Could not read file")
		return "", err
	}

	// Send start worker message to the master
	c := Command{0, "startWorker", f}
	msg, err := json.Marshal(c)
	if err != nil {
		fmt.Println("Could not marshal start worker message")
		return "", err
	}

	requester.Send(string(msg), 0)

	// Get reply which contains the port where the new worker runs
	reply, err := requester.Recv(0)
	if err != nil {
		fmt.Println("Did not recv worker port")
		return "", err
	}

	worker := new(ComputeResponse)

	err = json.Unmarshal([]byte(reply), worker)
	if err != nil {
		fmt.Println("Could not unmarshal response from compute", reply)
		return "", err
	}

	if worker.Status != 1 {
		err = errors.New("Worker not started")
		return "", err
	}

	return worker.Message, nil
}

func PrepareResponse(resp interface{}) []string {

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
