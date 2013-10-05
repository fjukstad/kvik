package kegg

import (
    "log"
    "github.com/fjukstad/gocache"
)

type Gene struct {
    Name string
    Definition string
    Orthology string
    Organism string
    Pathway []string
    Class []string
    Position string
    DBLinks map[string]string
    Structure string
    AASEQ Sequence
    NTSEQ Sequence
}

type Sequence struct {
    length int
    Sequence string
}

func getGene(id string) string {

    baseURL := "http://rest.kegg.jp/get/"
    url := baseURL + id

    
    response, err := gocache.Get(url)
    if err != nil{
        log.Panic("Cannot download from url:",err)
    }
    
    return "hei"
}

