package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/fjukstad/kegg"

	"code.google.com/p/gorest"
	"github.com/fjukstad/rpcman"
)

var rpc *rpcman.RPCMan

type RestService struct {

	// REST service details
	gorest.RestService `root:"/"
                        consumes:"application/json"
                        produces: "application/json" `

	geneExpression gorest.EndPoint `method:"GET"
                                    path:"/gene/{Id:string}"
                                    output:"[]float64"`

	avgDiff gorest.EndPoint `method:"GET"
                            path:"/gene/{Id:string}/avg"
                            output:"float64"`

	avgDiffs gorest.EndPoint `method:"GET"
                            path:"/genes/{...:string}/avg"
                            output:"string"`

	std gorest.EndPoint `method:"GET"
                            path:"/gene/{Id:string}/stddev"
                            output:"float64"`

	variance gorest.EndPoint `method:"GET"
                            path:"/gene/{Id:string}/vari"
                            output:"float64"`

	setScale gorest.EndPoint `method:"POST"
                             path:"/setscale/"
                             postdata:"string"`

	bg gorest.EndPoint `method:"GET"
                        path:"/gene/{GeneId:string}/{Exprs:string}/bg"
                        output:"string"`

	setSettings gorest.EndPoint `method:"POST" 
								path:"/setsettings/"
								postdata:"string"`

	getSettings gorest.EndPoint `method:"GET"
	path:"/getsettings/{Things:string}"
								output:"string"`

	// Dataset holding nowac data
	Dataset *Dataset

	// Settings
	Settings *Settings

	// RPC Server for performing statistics
	RPC *rpcman.RPCMan
}

type Settings struct {
	Smoking        bool
	HormoneTherapy bool
	Disable        bool
}

type Ex struct {
	Expression map[string]float64
}

func (serv RestService) GetSettings(Things string) string {

	j, err := json.Marshal(serv.Settings)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(j)

}

func (serv RestService) SetSettings(PostData string) {

	j := []byte(PostData)

	s := new(Settings)
	err := json.Unmarshal(j, &s)
	if err != nil {
		log.Println("Could not parse JSON settings", err)
	}

	serv.Settings.Disable = s.Disable
	serv.Settings.Smoking = s.Smoking
	serv.Settings.HormoneTherapy = s.HormoneTherapy

}

func (serv RestService) AvgDiffs(args ...string) string {

	if len(args) < 2 {
		return ""
	}

	// Last item in args is avg
	genes := args[0 : len(args)-1]
	genes = strings.Split(genes[0], " ")
	resp := make(map[string]float64, 0)

	for _, gene := range genes {
		resp[gene] = serv.AvgDiff(gene)
	}

	exp := Ex{resp}

	b, err := json.Marshal(exp)
	if err != nil {
		log.Println("Could not marshal reply in avgdiffs", err)
		return ""
	}

	return string(b)

}

func (serv RestService) SetScale(PostData string) {
	log.Print("setting scale to ", PostData)

	serv.Dataset.setScale(PostData)

	log.Print("---------------------------------------------")
}

func (dataset *Dataset) setScale(scale string) {
	// changing scale to the same as before
	if dataset.Scale == scale {
		return
	}

	log.Println("--------- old scale -------- ", dataset.Scale)
	dataset.Scale = scale
	log.Println("--------- new scale -------- ", dataset.Scale)

	var tmpIdExprs map[string][]float64
	var tmpGeneExprs map[string]map[string]*CaseCtrl

	tmpIdExprs = dataset.Exprs.IdExpression
	tmpGeneExprs = dataset.Exprs.GeneExpression

	dataset.Exprs.IdExpression = dataset.Exprs.DiffIdExpression
	dataset.Exprs.GeneExpression = dataset.Exprs.DiffGeneExpression

	dataset.Exprs.DiffIdExpression = tmpIdExprs
	dataset.Exprs.DiffGeneExpression = tmpGeneExprs

}

// convert
func log2(input []float64) []float64 {
	new_vals := make([]float64, len(input))

	for i, value := range input {
		v := math.Log2(value)
		new_vals[i] = v
	}

	return new_vals
}

func exp2(input []float64) []float64 {

	new_vals := make([]float64, len(input))

	for i, value := range input {
		v := math.Exp2(value)
		new_vals[i] = v
	}

	return new_vals
}

// Get gene expression for given gene
func (serv RestService) GeneExpression(Id string) []float64 {
	log.Print("Returning gene expression for gene ", Id)
	id := strings.Trim(Id, "hsa:")
	gene := kegg.GetGene(id)

	log.Print("hsa:", id, " ==> ", gene.Name)

	if gene.Name == "" {
		log.Println("Gene with id ", Id, " not found in database")

		// return slice with all zeros
		emptySlice := make([]float64, len(serv.Dataset.Exprs.Genes))
		return emptySlice
	}

	name := strings.Split(gene.Name, ", ")[0]

	var ret []float64
	// return difference between case & ctrl
	for _, cc := range serv.Dataset.Exprs.GeneExpression[name] {
		questions := serv.Dataset.Qs.Results[cc.Id]

		if serv.Settings.Disable {
			ret = append(ret, cc.Case-cc.Ctrl)
			continue
		}

		if serv.Settings.Smoking && questions.SmokingStatus == "Yes" {
			if serv.Settings.HormoneTherapy && questions.HormoneTherapy == "Yes" {
				ret = append(ret, cc.Case-cc.Ctrl)
				continue
			} else if !serv.Settings.HormoneTherapy && questions.HormoneTherapy == "No" {
				ret = append(ret, cc.Case-cc.Ctrl)
				continue
			}
		}

		if !serv.Settings.Smoking && questions.SmokingStatus == "No" {
			if !serv.Settings.HormoneTherapy && questions.HormoneTherapy == "No" {
				ret = append(ret, cc.Case-cc.Ctrl)
			} else if serv.Settings.HormoneTherapy && questions.HormoneTherapy == "Yes" {
				ret = append(ret, cc.Case-cc.Ctrl)
				continue
			}

		}
	}

	/*
	   if(serv.Context != nil){
	       serv.RB().ConnectionClose()
	   }
	*/

	return ret
}

// Get standard deviation for expression values of a given gene
func (serv RestService) Std(GeneId string) float64 {
	exprs := serv.GeneExpression(GeneId)

	if len(exprs) == 0 {
		log.Print("Expression values for gene ", GeneId, " not found")
		return 0
	}
	ret, err := serv.RPC.Call("std", exprs)
	//ret, err := serv.RPC.Call("add", 2, 5)
	if err != nil {
		log.Println("RPC FAILED", err)
		return 0
	}

	std, ok := ret.(float64)
	if !ok {
		log.Println("conversion to float64 went bad: ", ret)
		return 0
	}

	log.Println("Standard deviation for expression of gene ", GeneId, " is ", std)
	return std

}

// Get variation for expression values of a given gene
func (serv RestService) Variance(GeneId string) float64 {
	exprs := serv.GeneExpression(GeneId)

	if len(exprs) == 0 {
		log.Print("Expression values for gene ", GeneId, " not found")
		return 0
	}

	ret, err := serv.RPC.Call("var", exprs)
	if err != nil {
		log.Println("RPC FAILED", err)
		return 0
	}

	variance, ok := ret.(float64)
	if !ok {
		log.Println("conversion to float64 went bad: ", ret)
		return 0
	}

	log.Println("Variance for expression of gene ", GeneId, " is ", variance)
	return variance

}

func (serv RestService) AvgDiff(Id string) float64 {
	exprs := serv.GeneExpression(Id)

	if len(exprs) == 0 {
		log.Print("Expression values for gene ", Id, " not found")
		return 0
	}

	// avg := avg(exprs)

	ret, err := serv.RPC.Call("mean", exprs)
	if err != nil {
		log.Println("RPC FAILED", err)
	}
	avg, ok := ret.(float64)
	if !ok {
		log.Println("conversion to float64 went bad: ", ret)
	}

	log.Println("Average difference for gene ", Id, " is ", avg)

	/*
	   if(serv.Context != nil){
	       serv.RB().ConnectionClose()
	   }
	*/

	return avg
}

func avg(nums []float64) float64 {

	var total float64

	for _, num := range nums {
		total += num
	}

	return total / float64(len(nums))

}

// Find dataset id that has the given expression value
func expressionToId(dataset *Dataset, GeneId, Exprs string) string {

	id := strings.Trim(GeneId, "hsa:")
	gene := kegg.GetGene(id)

	if gene.Name == "" {
		log.Println("Gene with id ", GeneId, " not found in database")
		return ""
	}

	name := strings.Split(gene.Name, ", ")[0]

	dsId := ""

	exprsVal, err := strconv.ParseFloat(Exprs, 64)
	if err != nil {
		log.Println("could not convert ", Exprs, "to float")
		return ""
	}

	// return difference between case & ctrl
	for i, cc := range dataset.Exprs.GeneExpression[name] {
		ex := cc.Case - cc.Ctrl
		if ex == exprsVal {
			dsId = i
			break
		}

	}

	return dsId

}

func (serv RestService) Bg(GeneId, Exprs string) string {
	dsId := expressionToId(serv.Dataset, GeneId, Exprs)

	bg := serv.Dataset.Bg.IdInfo[dsId]

	b, err := json.Marshal(bg)
	if err != nil {
		log.Print("marshaling went bad: ", err)
		return ""
	}

	return string(b)
}

func Init(path string) *RestService {
	ds := NewDataset(path) //Dataset{} // := NewDataset(*path)

	ds.PrintDebugInfo()

	restService := new(RestService)
	restService.Dataset = &ds

	var err error

	// connect to statistics engine that will run statistics and that
	rpcaddr := []string{"tcp://localhost:5555", "tcp://localhost:5556"}
	//	"tcp://localhost:5557", "tcp://localhost:5558"} // "ipc:///tmp/datastore/0"
	restService.RPC, err = rpcman.Init(rpcaddr)

	// I know it's a global thing. sorry...
	rpc = restService.RPC

	if err != nil {
		log.Println("RPC error", err)
	}

	// Set up the settings
	restService.Settings = new(Settings)
	restService.Settings.Smoking = true
	restService.Settings.HormoneTherapy = true
	restService.Settings.Disable = true

	return restService
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

func main() {

	var path = flag.String("path", "/Users/bjorn/stallo/data", "path where data files are stored")
	var ip = flag.String("ip", "localhost", "ip to run on")
	var port = flag.String("port", ":8888", "port to run on")

	flag.Parse()

	restService := Init(*path)
	log.Print("Starting datastore at ", *ip, *port)

	http.HandleFunc("/com", CommandHandler)

	gorest.RegisterService(restService)
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(*port, nil)

}
