package main

import (
	"encoding/json"
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

type Response struct {
	Genes   []string
	Results []string
}

func InitRPC(address []string) (*rpcman.RPCMan, error) {
	rpc, err := rpcman.Init(address)

	if err != nil {
		return rpc, err
	}

	return rpc, nil

}

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

func RunCommand(command string) []string {
	output, err := rpc.Call("command", command)
	if err != nil {
		log.Println(err)
	}

	var result string
	t := reflect.TypeOf(output).String()

	if t == "string" {
		result = output.(string)
	} else {
		log.Println("ERROR, could not parse output from data-engine")
		return []string{""}
	}

	result = strings.Trim(result, "[1] ") // remove r thing
	result = strings.Trim(result, "\n")   // unwanted newlines
	log.Println(result)
	results := strings.Split(result, "  ")

	return results
}

func WriteResponse(w http.ResponseWriter, response *Response) {
	b, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(b)

}

func FoldChangeHandler(w http.ResponseWriter, r *http.Request) {
	list := strings.Split(r.URL.Path, "/fc/")[1]
	genes := strings.Split(list, "+")

	command := "fc(" + ListToVector(genes) + ")"

	results := RunCommand(command)

	response := new(Response)
	response.Genes = genes
	response.Results = results

	WriteResponse(w, response)

}

func PValueHandler(w http.ResponseWriter, r *http.Request) {

}

func GeneExpressionHandler(w http.ResponseWriter, r *http.Request) {

}

func ScaleHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {

	var ip = flag.String("ip", "localhost", "ip to run on")
	var port = flag.String("port", ":8888", "port to run on")

	flag.Parse()

	log.Print("Starting datastore at ", *ip, *port)

	var err error
	rpc, err = rpcman.Init([]string{"tcp://localhost:5555"})
	if err != nil {
		log.Panic(err)
	}

	http.HandleFunc("/com", CommandHandler)
	http.HandleFunc("/fc/", FoldChangeHandler)
	http.HandleFunc("/pval/", PValueHandler)
	http.HandleFunc("/exprs/", GeneExpressionHandler)
	http.HandleFunc("/scale/", ScaleHandler)

	err = http.ListenAndServe(*port, nil)
	if err != nil {
		log.Print(err)
	}

}
