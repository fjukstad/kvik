package utils

import (
	"log"
	"os"
	"path"
	"strings"
)

type Command struct {
	Type    int
	Command string
	File    []byte
}

type ComputeResponse struct {
	Status  int
	Message string
}

type WorkerRequest struct {
	Method string
	Args   interface{}
}

type WorkerResponse struct {
	Response interface{}
	Status   int
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
