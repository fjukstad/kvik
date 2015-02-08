package utils

import (
	"log"
	"os"
	"path"
	"strings"
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
	Output map[string]string
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
