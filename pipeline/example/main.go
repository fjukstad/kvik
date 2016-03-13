package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/fjukstad/kvik/pipeline"
	"github.com/fjukstad/kvik/r"
)

func main() {
	addr := "localhost:8181"
	//addr := "public.opencpu.org"
	username := ""
	password := ""

	var filename = flag.String("pipeline", "/Users/bjorn/Dropbox/go/src/github.com/fjukstad/kvik/pipeline/example/x-y/pipeline.yaml", "the pipeline description")

	flag.Parse()

	s := r.Server{addr, username, password}
	p, err := pipeline.ImportPipeline(*filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	p.RServer = &s

	runtime.GOMAXPROCS(runtime.NumCPU())

	_, err = p.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	p.Save()

	p.Print()

	re, err := p.Results("pdf")

	if err != nil {
		fmt.Println(err)
		return
	}

	ioutil.WriteFile("output.pdf", re, 0755)

	pkgs, err := s.InstalledPackages()
	if err != nil {
		fmt.Println(err)
		return
	}

	ioutil.WriteFile("pkgs.json", pkgs, 0755)

}
