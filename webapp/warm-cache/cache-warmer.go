package main

import (
	"fmt"

	"github.com/fjukstad/kegg"
)

func main() {

	//var host = flag.String("host", "localhost", "host where Kvik is running")

	//flag.Parse()

	//url := "http://" + *host + "/browser/"

	pws := kegg.GetAllHumanPathways()
	for _, p := range pws {

		//_ := "http://" + *host + "/browser/pathwaySelect=" + p

		fmt.Println(p)
	}

}
