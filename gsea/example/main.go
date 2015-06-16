package main

import "github.com/fjukstad/kvik/gsea"
import "fmt"

func main() {
	abstract, _ := gsea.Abstract("BIOCARTA_HER2_PATHWAY")
	fmt.Println(abstract)

	brief, _ := gsea.BriefDescription("PID_HES_HEY_PATHWAY")
	fmt.Println(brief)
}
