package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/fjukstad/r"
	"github.com/gorilla/mux"
)

var cache map[string]string
var logTime = "02-01-2006 15:04:05.000"

func CallHandler(w http.ResponseWriter, r *http.Request) {
	printTime()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	call := R.RCall{}
	err = json.Unmarshal(body, &call)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Call:", call.Package, call.Function, call.Arguments)

	cacheKey := call.Package + call.Function + call.Arguments

	if cache[cacheKey] != "" {
		printTime()
		fmt.Println("CACHE HIT")
		w.Write([]byte(cache[cacheKey]))
		return
	} else {
		printTime()
		fmt.Println("CACHE MISS")
	}

	s, err := R.Call(call.Package, call.Function, call.Arguments)
	if err != nil {
		fmt.Println(err)
		return
	}

	if cache == nil {
		cache = make(map[string]string, 0)
	}

	cache[cacheKey] = s.Key

	w.Write([]byte(s.Key))

}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	printTime()
	vars := mux.Vars(r)
	key := vars["key"]
	format := vars["format"]

	fmt.Println("Get: key:", key, "format", format)

	res, err := R.Get(key, format)

	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(res)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("fuck you"))
}

func printTime() {
	fmt.Print(time.Now().Format(logTime), " ")
}

func main() {

	R.Init("/tmp/go")

	r := mux.NewRouter()

	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/call", CallHandler)
	r.HandleFunc("/get/{key}/{format}", GetHandler)

	http.Handle("/", r)

	err := http.ListenAndServe(":8181", r)
	if err != nil {
		fmt.Println(err)
	}
}
