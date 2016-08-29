package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/fjukstad/kvik/r"
	"github.com/gorilla/mux"
)

func main() {
	port := flag.String("port", ":8181", "runs server on specified port")
	path := flag.String("dir", "/tmp/kvik", "tmp dir")
	flag.Parse()

	s, err := r.InitServer(10, *path)
	if err != nil {
		fmt.Println(err)
	}

	s.EnableCaching()

	router := mux.NewRouter()

	router.HandleFunc("/call", s.CallHandler)
	router.HandleFunc("/get/{key}/{format}", s.GetHandler)
	http.Handle("/", router)

	fmt.Println(http.ListenAndServe(*port, router))

}
