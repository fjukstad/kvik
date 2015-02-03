package main

import (
	"fmt"

	"github.com/fjukstad/kvik/genecards"
)

func main() {
	sum := genecards.Summary("BRCA1")
	fmt.Println(sum)
}
