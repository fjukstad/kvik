package main

import (
	"fmt"
	"time"

	"github.com/fjukstad/kvik/gopencpu"
)

func main() {

	g := gopencpu.GoOpenCPU{"http://docker0.bci.mcgill.ca:8004", "", ""}

	t0 := time.Now()
	a, err := g.Call("stats/R/rnorm", `{"n":100}`)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := a.GetResult(&g, "json")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)

	err = g.Plot("graphics/R/plot", `x=`+a.Key, "pdf", "output.pdf")
	if err != nil {
		fmt.Println(err)
	}

	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))

}
