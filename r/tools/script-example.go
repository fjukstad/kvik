package main

import (
	"fmt"

	"github.com/fjukstad/kvik/r"
)

func main() {

	err := r.StartServer(":8080", "scriptdir")
	fmt.Println(err)

}
