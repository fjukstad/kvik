package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"testing"
	//    "strings"
	//    "time"
	"flag"

	"code.google.com/p/gorest"
)

var ds Dataset
var restService *RestService

func BenchmarkInit(b *testing.B) {

	size := flag.Lookup("test.cpuprofile").Value.String()

	//disableLogging()
	path := "/Users/bjorn/stallo/src/src/nowac/datastore/data" + size
	// Set up rest service to be used later. reset timer afterwards
	restService = Init(path)
	gorest.RegisterService(restService)

	http.Handle("/", gorest.Handle())
	go http.ListenAndServe(":8888", nil)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		ds = NewDataset(path)
	}
}

var result []float64
var err error

func BenchmarkGeneExpression(b *testing.B) {
	//fmt.Println("")
	for n := 0; n < b.N; n++ {
		reply, err = get("http://localhost:8888/gene/hsa:23533")
		//fmt.Print(reply, "\r")
	}
}

var res float64

func BenchmarkAverageDiff(b *testing.B) {
	for n := 0; n < b.N; n++ {
		reply, err = get("http://localhost:8888/gene/hsa:23533/avg")
	}
}

func BenchmarkStd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		reply, err = get("http://localhost:8888/gene/hsa:23533/stddev")
	}
}

func BenchmarkVariance(b *testing.B) {
	for n := 0; n < b.N; n++ {
		reply, err = get("http://localhost:8888/gene/hsa:23533/vari")
	}
}

func BenchmarkSetScale(b *testing.B) {
	for n := 0; n < b.N; n++ {
		reply, err = get("http://localhost:8888/setscale/log")
	}
}

var exprs, bg string

func BenchmarkBackground(b *testing.B) {
	result = restService.GeneExpression("hsa:23533")
	b.ResetTimer()

	// get expression value from prev test
	exprs = strconv.FormatFloat(result[0], 'f', -1, 64)

	for n := 0; n < b.N; n++ {
		bg = restService.Bg("hsa:23533", exprs)
	}
}

func disableLogging() {
	log.SetOutput(ioutil.Discard)
}

var genes = [...]string{"hsa:1630", "hsa:836", "hsa:842", "hsa:836", "hsa:999", "hsa:1499", "hsa:51384 hsa:54361 hsa:7471 hsa:7472 hsa:7473 hsa2 hsa:7483 hsa:7484 hsa:80326 hsa:81029 hsa:89780", "hsa:1495 hsa:1496 hsa:29119", "hsa:1499", "hsa:3688", "hsa:1855 hsa:1856 hsa:1857", "hsa:10319 hsa:1282 hsa:1284 hsa:1285 hsa:1286 hsa:3911 hsa:3912 hsa:3913 hsa:3914 hsa:3915 hsa:3918", "hsa:3655 hsa:3673 hsa:3674 hsa:3675 hsa:3685", "hsa:5747", "hsa:5728", "hsa:11211 hsa:2535 hsa:7855 hsa:7976 hsa:8321 hsa:8322 hsa:8323 hsa:8324 hsa:8325 hsa:8326", "hsa:26060", "hsa:25 hsa:613", "hsa:1398 hsa:1399", "hsa:23624 hsa:867 hsa:868", "hsa:1398 hsa:1399", "hsa:23533 hsa:5290 hsa:5291 hsa:5293 hsa:5294 hsa:5295 hsa:5296 hsa:8503", "hsa:6776 hsa:6777", "hsa:3716", "hsa:1956", "hsa:2064", "hsa:5156 hsa:5159", "hsa:3480", "hsa:3815", "hsa:2322", "hsa:4233", "hsa:2260 hsa:2261 hsa:2263", "hsa:10342 hsa:4914 hsa:7170 hsa:7175", "hsa:7039", "hsa:1950", "hsa:3082", "hsa:5154 hsa:5155", "hsa:3479", "hsa:367", "hsa:3320 hsa:3326 hsa:7184", "hsa:367", "hsa:367", "hsa:2932", "hsa:8312 hsa:8313", "hsa:10297 hsa:324", "hsa:1147 hsa:3551 hsa:8517", "hsa:2475", "hsa:842", "hsa:572", "hsa:4193", "hsa:1027", "hsa:598", "hsa:6774", "hsa:6772", "hsa:2885", "hsa:6654 hsa:6655", "hsa:3265 hsa:3845 hsa:4893", "hsa:369 hsa:5894 hsa:673", "hsa:5604 hsa:5605", "hsa:83593", "hsa:11186", "hsa:4792", "hsa:4790 hsa:4791 hsa:5970", "hsa:6789", "hsa:595", "hsa:5979 hsa:8030 hsa:8031", "hsa:11186", "hsa:5335 hsa:5336", "hsa:5900", "hsa:5898 hsa:5899", "hsa:387 hsa:5879 hsa:5880 hsa:5881", "hsa:10928", "hsa:5337", "hsa:5599 hsa:5601 hsa:5602", "hsa:5879 hsa:5880 hsa:5881 hsa:998", "hsa:3725", "hsa:2353", "hsa:4609", "hsa:7157", "hsa:51176 hsa:6932 hsa:6934 hsa:83439", "hsa:5578 hsa:5579 hsa:5582", "hsa:6469", "hsa:5727", "hsa:6608", "hsa:27148", "hsa:51684", "hsa:2271", "hsa:112398 hsa:112399 hsa:54583", "hsa:6923", "hsa:6921", "hsa:7428", "hsa:8453", "hsa:9978", "hsa:2034 hsa:3091", "hsa:405 hsa:9915", "hsa:1387 hsa:2033", "hsa:54205", "hsa:675", "hsa:5888", "hsa:4292"}

var reply string

/*
func BenchmarkManyGenes(b *testing.B){
    disableLogging()
    path := "/Users/bjorn/stallo/src/src/nowac/datastore/data"
    // Set up rest service to be used later. reset timer afterwards
    restService = Init(path)


    fmt.Println("HERE WE ARE AGAIN")
    gorest.RegisterService(restService)

    http.Handle("/", gorest.Handle())
    go http.ListenAndServe(":8888", nil)


    startup := time.Now()

    b.ResetTimer()

    numRuns := 3000
    for n := 0; n < b.N; n++ {
        for j := 0; j < numRuns; j++ {
            sta := time.Now()
            for _, gene := range(genes){
                g := strings.Split(gene, " ")[0]
                reply = get("http://localhost:8888/gene/"+g)
                reply = get("http://localhost:8888/gene/"+g+"/avg")
            }
          elapsed := time.Since(sta)
          fmt.Println("Done with genes in ",elapsed.Seconds(),"seconds. (",
                        time.Since(startup)," total runtime)")

         log.Println(" --------------------",j,"-----------------")
       }
    }
}
*/

func get(URL string) (string, error) {
	resp, err := http.Get(URL)

	if err != nil {
		fmt.Println("HTTP GET ERROR:", err)
		return "ERROR", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Print("Reading response body went bad. ", err)
		return "ERROR", err
	}

	return string(body), nil

}

/*
func BenchmarkDatasetSize1x(b *testing.B) {
    path := "/Users/bjorn/stallo/src/src/nowac/datastore/data"

    for n := 0; n < b.N; n++ {
        ds = NewDataset(path)
    }
}


func BenchmarkDatasetSize2x(b *testing.B) {
    size := "2x"
    path := "/Users/bjorn/stallo/src/src/nowac/datastore/data-"+size

    for n := 0; n < b.N; n++ {
        ds = NewDataset(path)
    }
}


func BenchmarkDatasetSize5x(b *testing.B) {
    size := "5x"
    path := "/Users/bjorn/stallo/src/src/nowac/datastore/data-"+size

    for n := 0; n < b.N; n++ {
        ds = NewDataset(path)
    }
}

func BenchmarkDatasetSize10x(b *testing.B) {
    size := "10x"
    path := "/Users/bjorn/stallo/src/src/nowac/datastore/data-"+size

    for n := 0; n < b.N; n++ {
        ds = NewDataset(path)
    }
}

func BenchmarkDatasetSize20x(b *testing.B) {
    size := "20x"
    path := "/Users/bjorn/stallo/src/src/nowac/datastore/data-"+size

    for n := 0; n < b.N; n++ {
        ds = NewDataset(path)
    }
}

func BenchmarkDatasetSize40x(b *testing.B) {
    size := "40x"
    path := "/Users/bjorn/stallo/src/src/nowac/datastore/data-"+size

    for n := 0; n < b.N; n++ {
        ds = NewDataset(path)
    }
}
*/
