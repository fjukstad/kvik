package main

import (
	"fmt"

	"github.com/fjukstad/kvik/r"
)

func main() {

	port := ":8181"
	err := r.StartServer(port, "/tmp/kvik")
	fmt.Println(err)
}
