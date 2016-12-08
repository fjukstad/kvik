package main

import (
	"fmt"

	"github.com/fjukstad/kvik/eutils"
)

func main() {
	r, err := eutils.Search([]string{"BRCA1", "BRCA2"})
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := eutils.Summary(r)

	for _, paper := range res.DocumentSummarySet.DocumentSummary {
		fmt.Println(paper.Title)
		fmt.Println(paper.ViewCount)
	}
}
