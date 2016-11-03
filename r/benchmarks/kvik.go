package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/fjukstad/kvik/r"
)

func main() {
	t0 := time.Now()

	s := r.Server{"localhost:80", "", ""}
	var out string
	var err error
	out, err = s.Call("stats", "rnorm", "n=100")

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := s.Get(out, "json")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)

	plot, err := s.Call("graphics", "plot", "x="+out)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := s.Get(plot, "pdf")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile("output.pdf", result, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Generated plot and saved it as output.pdf")

	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))

}
