package main

import (
	"flag"
	"fmt"

	"github.com/fjukstad/kvik/r"
)

func main() {

	port := flag.String("port", ":8181", "runs server on specified port")
	path := flag.String("dir", "/tmp/kvik", "tmp dir")
	flag.Parse()
	err := r.StartServer(*port, *path)
	fmt.Println(err)
}
