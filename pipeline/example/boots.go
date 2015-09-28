package main

import (
	"fmt"
	"strconv"

	"github.com/fjukstad/kvik/kompute"
	"github.com/fjukstad/kvik/pipeline"
)

func main() {
	addr := "192.168.99.100:8004"
	//addr := "public.opencpu.org:80"
	username := "user"
	password := "password"

	k := kompute.NewKompute(addr, username, password)

	p := pipeline.NewPipeline("boots", k)

	// --------------- LOAD DATA ---------------- //
	name := "loaddata"
	function := "syntheticdata"
	pkg := "github.com/fjukstad/boots"
	argnames := []string{"nsamples", "class", "noisevars"}
	args := []string{
		"1050",
		"T",
		"9000",
	}
	s := pipeline.NewStage(name, function, pkg, argnames, args)
	p.AddStage(s)

	name = "response"
	function = "responses"
	pkg = "github.com/fjukstad/boots"
	argnames = []string{"dataset"}
	args = []string{
		"from:loaddata",
	}

	s = pipeline.NewStage(name, function, pkg, argnames, args)
	p.AddStage(s)

	name = "predictors"
	function = "predictors"
	pkg = "github.com/fjukstad/boots"
	argnames = []string{"dataset"}
	args = []string{"from:loaddata"}

	s = pipeline.NewStage(name, function, pkg, argnames, args)
	p.AddStage(s)

	numBoots := 3

	for i := 0; i < numBoots; i++ {
		name = "boots-" + strconv.Itoa(i)
		function = "boots"
		pkg = "github.com/fjukstad/boots"
		argnames = []string{"X", "Y"}
		args = []string{"from:predictors", "from:response"}

		s = pipeline.NewStage(name, function, pkg, argnames, args)
		p.AddStage(s)
	}

	p.Run()

	p.Print()
	fmt.Println("done...")

	p.Save()

}
