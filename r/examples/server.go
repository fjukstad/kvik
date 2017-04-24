package main

import (
	"flag"
	"fmt"

	"github.com/fjukstad/kvik/r"
)

func main() {
	port := flag.String("port", ":8181", "runs server on specified port")
	path := flag.String("dir", "/tmp/kvik", "tmp dir")
	cache := flag.Bool("cache", false, "enable caching of results")
	flag.Parse()

	s, err := r.InitServer(10, *path)
	if err != nil {
		fmt.Println(err)
	}

	if *cache {
		s.EnableCaching()
	}
	fmt.Println(s.Start(*port))
}
