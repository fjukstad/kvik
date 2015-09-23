package main

import (
	"fmt"
	"runtime"

	"github.com/fjukstad/kvik/kompute"
	"github.com/fjukstad/kvik/pipeline"
)

func main() {
	addr := "192.168.99.100:8004"
	username := "user"
	password := "password"

	k := kompute.NewKompute(addr, username, password)
	p, err := pipeline.ImportPipeline("pipeline.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	p.Kompute = k

	runtime.GOMAXPROCS(runtime.NumCPU())

	res, err := p.Run()
	if err != nil {
		fmt.Println(err)
	}

	for _, r := range res {
		fmt.Println(r.Key)
	}

	p.Save()

	p.Print()

	_, err = p.Results("png")

	if err != nil {
		fmt.Println(err)
	}

}
