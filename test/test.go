package main

import (
    "nowac/kegg"
    "log"
)

func main(){

    geneId := "hsa:10458"
    gene := kegg.GetGene(geneId);

    log.Print("get:",geneId)
    log.Print(gene)
   
}
