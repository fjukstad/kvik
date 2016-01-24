package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/fjukstad/kvik/kompute"
	"github.com/fjukstad/kvik/pipeline"
)

func main() {
	//addr := "192.168.99.100:8004"
	addr := "public.opencpu.org"
	username := "user"
	password := "password"

	var filename = flag.String("pipeline", "pipeline.yaml", "the pipeline description")

	flag.Parse()

	k := kompute.NewKompute(addr, username, password)
	p, err := pipeline.ImportPipeline(*filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	p.Kompute = k

	runtime.GOMAXPROCS(runtime.NumCPU())

	_, err = p.Run()
	if err != nil {
		fmt.Println(err)
	}

	p.Save()

	p.Print()

	_, err = p.Results("png")

	if err != nil {
		fmt.Println(err)
	}

}
