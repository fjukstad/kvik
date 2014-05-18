package rpcman

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	zmq "github.com/alecthomas/gozmq"
)

type RPCMan struct {
	Context  *zmq.Context
	ServAddr string
	Socket   *zmq.Socket
	mutex    *sync.Mutex
}

type Request struct {
	Method string
	Args   interface{}
}

type Response struct {
	Response interface{}
	Status   int
}

func Init(addr string) (*RPCMan, error) {
	context, _ := zmq.NewContext()
	major, minor, patch := zmq.Version()

	socket, err := context.NewSocket(zmq.REQ)
	if err != nil {
		log.Println("could not set up socket to rpcman", err)
		return nil, err
	}
	err = socket.Connect(addr)
	if err != nil {
		log.Println("Could not connect to socket", err)
		return nil, err
	}

	// set timeout to 1s
	duration, _ := time.ParseDuration("1s")
	socket.SetRcvTimeout(duration)

	mutex := &sync.Mutex{}
	//cond := sync.NewCond(sync.Mutex{})

	rpc := RPCMan{
		Context:  context,
		ServAddr: addr,
		Socket:   socket,
		mutex:    mutex}

	return &rpc, nil

}

func (rpc RPCMan) Close() {
	rpc.Context.Close()
}

func (rpc RPCMan) Call(method string, args ...interface{}) (interface{}, error) {

	// exclusive access to socket
	//rpc.mutex.Lock()
	//defer rpc.mutex.Unlock()

	msg := Request{method, args}

	enc, err := json.Marshal(msg)
	if err != nil {
		log.Println("json error", err)
		return -1, err
	}

	rpc.Socket.Send(enc, 0)

	//rpc.Socket.Send(enc, 0)
	resp := new(Response)

	//reply, err := rpc.Socket.Recv(0)
	reply, err := rpc.Socket.Recv(0)
	if err != nil {
		log.Println("RPC recv failed:", err, method)
		return -1, err
	}

	err = json.Unmarshal(reply, resp)

	if err != nil {
		log.Println("json error", err)
		return -1, err
	}

	if resp.Status != 0 {
		return -1, errors.New("RPC Method not found:" + method)
	}

	// socket.Close()

	return resp.Response, nil

}
