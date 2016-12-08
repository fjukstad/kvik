package main

import (
	"fmt"

	"github.com/fjukstad/kvik/genenames"
)

func main() {
	symbol := "BRCA1"
	d, err := genenames.GetDoc(symbol)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(d)
}
