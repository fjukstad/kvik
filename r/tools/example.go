package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/fjukstad/kvik/r"
)

func main() {

	r.Init("/tmp/go", "")

	t0 := time.Now()

	pkg := "stats"
	fun := "rnorm"
	args := "n=100"

	s, err := r.Call(pkg, fun, args)
	if err != nil {
		fmt.Println(err)
		return
	}
	//	fmt.Println(s.Output)

	result, err := r.Get(s.Key, "json")
	if err != nil {
		fmt.Println("could not get csv")
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile("output.json", result, 0755)
	if err != nil {
		fmt.Println("could not write csv")
		fmt.Println(err)
		return
	}

	t1 := time.Now()
	fmt.Printf("The R (internal) call took %v to run.\n", t1.Sub(t0))

	b, err := r.Call(pkg, fun, args)
	if err != nil {
		fmt.Println("shit")
		fmt.Println(err)
		return
	}
	fmt.Println(b.Output)

	pkg = "graphics"
	fun = "plot"
	args = "x=" + s.Key + ",y=" + b.Key

	plot, err := r.Call(pkg, fun, args)
	if err != nil {
		fmt.Println("plot call failed")
		fmt.Println(err)
		return
	}

	result, err = r.Get(plot.Key, "pdf")
	if err != nil {
		fmt.Println("Could not generate pdf")
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile("output.pdf", result, 0755)
	if err != nil {
		fmt.Println("could not save pdf")
		fmt.Println(err)
		return
	}

	pkgs, err := r.InstalledPackages()
	fmt.Println(string(pkgs))
}
