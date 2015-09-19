package main

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/fjukstad/kvik/kompute"
	"github.com/fjukstad/kvik/pipeline"
)

func main() {
	addr := "192.168.99.100:8004"
	username := "user"
	password := "password"

	k := kompute.NewKompute(addr, username, password)

	p := pipeline.NewPipeline("pied piper", k)

	numStages := 50

	for i := 0; i < numStages; i++ {
		name := "stage-" + strconv.Itoa(i)
		function := "+"
		pkg := "base"
		argnames := []string{"e1", "e2"}
		args := []string{"200", "400"}
		s := pipeline.NewStage(name, function, pkg, argnames, args)
		p.AddStage(s)
	}

	name := "final stage"
	function := "+"
	pkg := "base"
	argnames := []string{"e1", "e2"}
	args := []string{"from:stage-0", "from:stage-1"}
	s := pipeline.NewStage(name, function, pkg, argnames, args)
	p.AddStage(s)

	runtime.GOMAXPROCS(runtime.NumCPU())

	res, err := p.Run()
	if err != nil {
		fmt.Println(err)
	}

	for _, r := range res {
		fmt.Println(r.Key)
	}

	p.Save()

	fmt.Println("done")
}
