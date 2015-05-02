package main

import (
	"fmt"

	"github.com/fjukstad/kvik/kompute"
)

func main() {

	k := kompute.Kompute{"http://opencpu:8004/", ""}

	a, err := k.Rpc("stats/R/rnorm", `{"n":10, "mean": 10, "sd":10}`, "json")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(a)

	b, err := k.Rpc("stats/R/rnorm", `{"n":10, "mean": 10, "sd":10}`, "csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(b)

	err = k.Plot("graphics/R/hist", `{"x": [1,2,5,1,20,12,11,5,4,6,10]}`,
		"png", "histogram.png")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("histogram stored in histogram.png")

}
