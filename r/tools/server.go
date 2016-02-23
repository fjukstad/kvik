package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/fjukstad/kvik/r"
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
	w.Write([]byte("Nothing to see here. Try /call, /get/{key}/{format} or /install"))
}

func InstallHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	defer file.Close()

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	fileLoc := "/tmp/file"

	out, err := os.Create(fileLoc)
	if err != nil {
		fmt.Fprintf(w, "Failed to open the file for writing")
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	printTime()
	fmt.Println("Installing package", header.Filename)

	output, err := R.InstallPackageFromSource(fileLoc)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	fmt.Fprintf(w, "%s", output)

}

func printTime() {
	fmt.Print(time.Now().Format(logTime), " ")
}

func main() {

	var packages = flag.String("packages", "", "packages you want installed, prereqs. Comma separated and surrounded by ''. E.g.: \"'dplyr', 'ggplot2'\" ")

	flag.Parse()

	printTime()
	fmt.Println("Installing packages", *packages)

	err := R.Init("/tmp/go", *packages)
	if err != nil {
		fmt.Println(err)
		return
	}

	r := mux.NewRouter()

	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/call", CallHandler)
	r.HandleFunc("/get/{key}/{format}", GetHandler)
	r.HandleFunc("/install", InstallHandler)

	http.Handle("/", r)

	printTime()
	port := ":8181"
	fmt.Println("Starting server at ", port)

	err = http.ListenAndServe(port, r)
	if err != nil {
		fmt.Println(err)
	}
}
