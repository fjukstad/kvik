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

	"github.com/fjukstad/kegg"
	"github.com/fjukstad/rpcman"
)

var rpc *rpcman.RPCMan

type Response struct {
	Result map[string]string
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

	command, err := url.QueryUnescape(strings.Split(r.URL.RawQuery, "=")[1])
	if err != nil {
		log.Panic("Could not unescape query:", err)
	}

	command = replaceOctets(command)

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
	} else if t == "float64" {
		res := output.(float64)
		return []string{strconv.FormatFloat(res, 'f', 5, 64)}
	} else {
		log.Println("ERROR, could not parse output from data-engine")
		return []string{""}
	}

	result = strings.Trim(result, "[1] ") // remove r thing
	result = strings.Trim(result, "\n")   // unwanted newlines
	result = strings.TrimLeft(result, " ")
	results := strings.Split(result, " ")

	var res []string
	for _, r := range results {
		if r != "" {
			res = append(res, r)
		}
	}

	return res
}

func WriteResponse(w http.ResponseWriter, genes []string, results []string) {
	response := new(Response)

	response.Result = make(map[string]string, len(genes))

	log.Println(results, genes)
	for i, gene := range genes {
		if results[i] != "NA" {
			response.Result[gene] = results[i]
		}

	}

	b, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(b)

}

func parseGenes(path string, separator string) ([]string, []string) {
	list := strings.Split(path, "/"+separator+"/")[1]
	geneIds := strings.Split(list, "+")

	var ids []string

	var genes []string

	for _, id := range geneIds {
		// Split on colon, if we dont have a colon we expect the string to be a
		// gene name and we can just return the list of these genes.
		gene := strings.Split(id, ":")
		if len(gene) < 2 {
			return geneIds, geneIds
		}
		id := gene[1]
		g := kegg.GetGene(id)
		names := strings.Split(g.Name, " ")
		for _, n := range names {
			name := strings.TrimRight(n, ",")
			genes = append(genes, name)
			ids = append(ids, id)
		}
	}
	return genes, ids

}

func generateCommand(name string, genes []string) string {
	command := name + "(" + ListToVector(genes) + ")"
	return command
}

func FoldChangeHandler(w http.ResponseWriter, r *http.Request) {
	com := "fc"
	genes, ids := parseGenes(r.URL.Path, com)
	command := generateCommand(com, genes)
	results := RunCommand(command)

	WriteResponse(w, ids, results)
}

func PValueHandler(w http.ResponseWriter, r *http.Request) {
	com := "pvalues"
	genes, ids := parseGenes(r.URL.Path, com)
	command := generateCommand(com, genes)
	results := RunCommand(command)
	WriteResponse(w, ids, results)
}

func GeneExpressionHandler(w http.ResponseWriter, r *http.Request) {
	com := "exprs"
	genes, _ := parseGenes(r.URL.Path, com)
	command := generateCommand(com, genes)
	results := RunCommand(command)
	log.Println("GOT EXPRESSION", genes, results)
	WriteExprsResponse(w, results)
}

func ScaleHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusInternalServerError)
}

type ExprsResponse struct {
	Exprs []string
}

func WriteExprsResponse(w http.ResponseWriter, exprs []string) {
	response := new(ExprsResponse)
	response.Exprs = exprs

	b, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(b)

}

type GeneResponse struct {
	Genes []string
}

func WriteGeneResponse(w http.ResponseWriter, genes []string) {
	response := new(GeneResponse)
	response.Genes = genes

	b, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(b)

}

func GeneHandler(w http.ResponseWriter, r *http.Request) {
	com := "genes()"
	results := RunCommand(com)
	WriteGeneResponse(w, results)
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
	http.HandleFunc("/pvalues/", PValueHandler)
	http.HandleFunc("/exprs/", GeneExpressionHandler)
	http.HandleFunc("/scale/", ScaleHandler)
	http.HandleFunc("/genes/", GeneHandler)

	err = http.ListenAndServe(*port, nil)
	if err != nil {
		log.Print(err)
	}

}
