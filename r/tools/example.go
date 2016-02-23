package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/fjukstad/r"
)

func main() {

	R.Init("/tmp/go")

	t0 := time.Now()

	pkg := "stats"
	fun := "rnorm"
	args := "n=10000000"

	s, err := R.Call(pkg, fun, args)
	if err != nil {
		fmt.Println(err)
		return
	}
	//	fmt.Println(s.Output)

	result, err = R.Get(s.Key, "json")
	if err != nil {
		fmt.Println("could not get csv")
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile("/Users/bjorn/Dropbox/go/src/github.com/fjukstad/r/example/output.json", result, 0755)
	if err != nil {
		fmt.Println("could not write csv")
		fmt.Println(err)
		return
	}

	t1 := time.Now()
	fmt.Printf("The R (internal) call took %v to run.\n", t1.Sub(t0))

	b, err := R.Call(pkg, fun, args)
	if err != nil {
		fmt.Println("shit")
		fmt.Println(err)
		return
	}
	fmt.Println(b.Output)

	pkg = "graphics"
	fun = "plot"
	args = "x=" + s.Key + ",y=" + b.Key

	plot, err := R.Call(pkg, fun, args)
	if err != nil {
		fmt.Println("plot call failed")
		fmt.Println(err)
		return
	}

	result, err = R.Get(plot.Key, "pdf")
	if err != nil {
		fmt.Println("Could not generate pdf")
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile("/Users/bjorn/Dropbox/go/src/github.com/fjukstad/r/example/output.pdf", result, 0755)
	if err != nil {
		fmt.Println("could not save pdf")
		fmt.Println(err)
		return
	}
}
