package main

import (
	"fmt"

	"github.com/fjukstad/kvik/gopencpu"
)

func main() {

	g := gopencpu.GoOpenCPU{"http://public.opencpu.org", "", ""}

	a, err := g.Rpc("stats/R/rnorm", `{"n":10, "mean": 10, "sd":10}`, "json")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(a)

	b, err := g.Rpc("stats/R/rnorm", `{"n":10, "mean": 10, "sd":10}`, "csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(b)

	err = g.Plot("graphics/R/hist", `{"x": [1,2,5,1,20,12,11,5,4,6,10]}`,
		"png", "histogram.png")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("histogram stored in histogram.png")

}
