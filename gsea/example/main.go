package main

import "github.com/fjukstad/kvik/gsea"
import "fmt"

func main() {
	a, err := gsea.Abstract("BIOCARTA_HER2_PATHWAY")

	fmt.Println(a, err)
}
