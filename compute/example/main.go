package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/fjukstad/kvik/utils"
	zmq "github.com/pebbe/zmq4"
)

func main() {

	var addr = flag.String("addr",
		"tcp://localhost:8888",
		"address to connect to")

	flag.Parse()

	requester, _ := zmq.NewSocket(zmq.REQ)
	defer requester.Close()
	err := requester.Connect(*addr)
	if err != nil {
		log.Panic(err)
	}

	f, err := ioutil.ReadFile("../worker/script.r")
	if err != nil {
		log.Panic(err)
	}
	sendCommand("startWorker", requester, f)
	sendCommand("killall", requester, nil)
}

func sendCommand(cmd string, requester *zmq.Socket, file []byte) {
	c := utils.Command{0, cmd, file}

	msg, err := json.Marshal(c)
	if err != nil {
		log.Panic(err)
	}

	requester.Send(string(msg), 0)

	reply, err := requester.Recv(0)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("reply:", reply)

}
