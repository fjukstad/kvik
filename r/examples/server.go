package main

import (
	"fmt"
	"net/http"

	"github.com/fjukstad/kvik/r"
	"github.com/gorilla/mux"
)

func main() {
	s, err := r.InitServer(4, "/Users/bjorn/Dropbox/go/src/github.com/fjukstad/kvik/r/tmp/kvikr")
	if err != nil {
		fmt.Println(err)
	}

	s.EnableCaching()

	router := mux.NewRouter()

	router.HandleFunc("/call", s.CallHandler)
	router.HandleFunc("/get/{key}/{format}", s.GetHandler)
	http.Handle("/", router)

	port := ":8181"

	fmt.Println(http.ListenAndServe(port, router))

}
