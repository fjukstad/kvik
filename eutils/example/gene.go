package main

import "fmt"
import "github.com/fjukstad/kvik/eutils"

func main() {
	ds, err := eutils.GeneSummary("627")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ds)
}
