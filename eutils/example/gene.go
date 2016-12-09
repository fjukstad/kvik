package main

import "fmt"
import "github.com/fjukstad/kvik/eutils"

func main() {
	ds, err := eutils.GetDocumentSummary("627")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ds)
}
