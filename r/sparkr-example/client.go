package main

import (
	"fmt"
	"io/ioutil"

	"github.com/fjukstad/kvik/r"
)

func main() {

	s := r.Server{":8181", "", ""}
	out, err := s.Upload("script.R", "script")
	fmt.Println("Upload return:", err, out)

	ioutil.WriteFile("output.zip", out, 0755)

}
