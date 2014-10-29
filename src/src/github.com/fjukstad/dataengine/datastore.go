package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/fjukstad/rpcman"
)

var rpc *rpcman.RPCMan

func replaceOctets(s string) string {
	res := strings.Replace(s, "%2B", "+", -1)
	res = strings.Replace(res, "%3D", "=", -1)
	res = strings.Replace(res, "%2C", ",", -1)
	res = strings.Replace(res, "%2F", "/", -1)
	return res
}

func CommandHandler(w http.ResponseWriter, r *http.Request) {
	// Cross origin nonsense
	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Println(r.URL.RawQuery)
	command, err := url.QueryUnescape(strings.Split(r.URL.RawQuery, "=")[1])
	if err != nil {
		log.Panic("Could not unescape query:", err)
	}

	command = replaceOctets(command)

	log.Println("Recieved command", command)

	// send command to rpcman
	output, err := rpc.Call("command", command)
	if err != nil {
		log.Println("RPC FAILED", err)
	}

	var response string
	t := reflect.TypeOf(output).String()

	if t == "float64" {
		res := output.(float64)
		response = strconv.FormatFloat(res, 'f', 3, 64)
	}
	if t == "string" {
		response = output.(string)
	}

	log.Println(response)

	_, err = io.WriteString(w, response)
	if err != nil {
		log.Panic("Error writing to response", err)
	}
}

func main() {

	var ip = flag.String("ip", "localhost", "ip to run on")
	var port = flag.String("port", ":8888", "port to run on")

	flag.Parse()

	log.Print("Starting datastore at ", *ip, *port)

	http.HandleFunc("/com", CommandHandler)

	http.ListenAndServe(*port, nil)

}
