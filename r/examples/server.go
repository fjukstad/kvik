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

	s, err := r.InitServer(10, *path)
	if err != nil {
		fmt.Println(err)
	}

	s.EnableCaching()
	fmt.Println(s.Start(*port))
}
