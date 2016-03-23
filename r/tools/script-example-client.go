package main

import (
	"fmt"

	"github.com/fjukstad/kvik/r"
)

func main() {

	s := r.Server{":8080", "", ""}
	err := s.Upload("script.R", "script")
	fmt.Println("Upload return:", err)

}
