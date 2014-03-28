package main

import (
    "log"
    "os"
    "encoding/csv"
    "io"
    "strconv"
    "path"
    "strings" 

)

type Dataset struct {
    Bg Background
    Exprs Expression

    Scale string
    DiffExprs Expression
}

type Background struct {

    IdInfo map[string] Info    
}

type Info struct {
    Lab int
    Id string
    Ctrl string
    Et  int
    En int
    Stage int
}

type Expression struct {
    Genes []string
    IdExpression map[string] [] float64
    GeneExpression map[string] map[string]*CaseCtrl
}

type CaseCtrl struct { 
    Case float64
    Ctrl float64
} 

func NewDataset( path string ) Dataset {

    exprsFilename := path + "/exprs.csv"
    bgFilename := path + "/background.csv"

    bg, err := generateBackgroundDataset(bgFilename)
    if err != nil {
        log.Print(err)
    }
    
    exprs, err := generateExpressionDataset(exprsFilename)
    if err != nil {
        log.Print(err) 
    }
    
    diffexprs := Expression{}
    // Init data set with background and expression data.  Set scale to absolute
    ds := Dataset{bg, exprs, "abs", diffexprs}

    return ds
    
}

func generateExpressionDataset(filename string) (Expression, error) { 

    var IdExpression map[string][]float64
    var GeneExpression map[string]map[string]*CaseCtrl
     
    exprs := Expression{}

    exprsfile, err := os.Open(filename)
    if err != nil {
        return exprs, err
    }  
    defer exprsfile.Close() 

    reader := csv.NewReader(exprsfile) 
    firstRow := true

    for {

        record, err := reader.Read() 
        
        if err == io.EOF{
            break
        } else if err != nil {
            log.Panic(err) 
        }
        if firstRow {
            exprs.Genes = probesToGenes(record, filename) 
            

            // the lengths here are maybe a bit off?
            IdExpression = make(map[string] []float64, len(record)-1)
            GeneExpression = make(map[string] map[string]*CaseCtrl, len(record)-1)
            for i, _ := range(exprs.Genes) { 
                GeneExpression[exprs.Genes[i]] = make(map[string] *CaseCtrl, len(record)) 
            } 
            firstRow = false 
        } else { 

            id := record[0]
            expression := toFloats(record[1:])
            
            // store an id to expression mapping. 
            IdExpression[id] = expression
            
            // string contains _ means it is a ctrl
            ctrl := strings.Contains(id, "_") 

            // Store the expression value for this specific gene and id
            // combination
            for i, _ := range(expression){
                
                gene := exprs.Genes[i]

                // Store id on same nuber without trailing _1 _2 etc
                commonId := strings.Split(id, "_")[0]

                exprsvals := GeneExpression[gene][commonId]

                if(exprsvals == nil) {
                    GeneExpression[gene][commonId] = new(CaseCtrl)
                } 

                if(ctrl){ 
                    GeneExpression[gene][commonId].Ctrl = expression[i]
                } else { 
                    GeneExpression[gene][commonId].Case = expression[i]
                } 
            }
        }
    }

    exprs.IdExpression = IdExpression
    exprs.GeneExpression = GeneExpression

    return exprs, nil
}

func probesToGenes(probes [] string, filename string) [] string {
    
    // assumes that filename has an extension
    path := path.Dir(filename) + "/"

    p2gfilename := path + "probe2gene.csv"


    p2g, err := getProbeToGeneMapping(p2gfilename) 

    if err != nil {
        log.Panic("Could not read probe 2 gene mappings! ", err)
    } 


    genes := make([] string, len(probes))

    for i , probeId := range(probes) {
        genes[i] = p2g[probeId]
    }


    return genes

}

func getProbeToGeneMapping(filename string) (map[string] string, error) { 
 
    // Note the 50 000 below. We'll have less than 50 000 genes, but if this
    // fails, then the whole thing breaks apart. 

    p2g := make(map[string] string, 50000)

    p2gfile, err := os.Open(filename)
    if err != nil {
        return p2g, err
    }  
    defer p2gfile.Close() 

    reader := csv.NewReader(p2gfile) 
    firstRow := true

    for {

        record, err := reader.Read() 
        
        if err == io.EOF{
            break
        } else if err != nil {
            log.Panic(err) 
        }
        if firstRow {
            firstRow = false
        }
        
        probeId := record[1]
        geneId := record[2]

        p2g[probeId] = geneId

    }


    return p2g, err

} 


func (ds Dataset) PrintDebugInfo() {
    exprs := ds.Exprs

    log.Print("Generated dataset with ", len(exprs.IdExpression["900229_1"]),
                " genes and ",len(exprs.GeneExpression[exprs.Genes[0]]),
                " case/ctrl pairs")

}


func generateBackgroundDataset(filename string) (Background, error) { 
    
    bg := Background{}

    bgfile, err := os.Open(filename)
    if err != nil {
        return bg, err
    }
    defer bgfile.Close() 

    
    reader := csv.NewReader(bgfile) 
    firstRow := true

    var idinfo map[string] Info

    for {

        record, err := reader.Read() 
        
        if err == io.EOF{
            break
        } else if err != nil {
            log.Panic(err) 
        }
        if firstRow {
            
            idinfo = make(map[string] Info, len(record)) 

            firstRow = false 
        } else { 
            var lab, et, en, stage int

            lab, err := strconv.Atoi(record[0])
            if err != nil {
                log.Print("Error parsing ") 
                lab = 0
            }

            id := record[1]
            casectrl := record[2]
            et, err = strconv.Atoi(record[3])
            if err != nil {
                et = -1
            }

            en, err = strconv.Atoi(record[4])
            if err != nil {
                en = -1
            }

            stage, err = strconv.Atoi(record[5])
            if err != nil {
                stage = -1
            }
            
            info := Info{
                Lab: lab,
                Id: id,
                Ctrl: casectrl,
                Et: et,
                En: en,
                Stage: stage,
            }
            idinfo[id] = info

        }
    }

    bg.IdInfo = idinfo


    return bg, nil
    
    
}



func toFloats(input []string) []float64{
    output := make([]float64, len(input))

    var err error

    for i, _ := range(input){
        output[i], err = strconv.ParseFloat(input[i], 32)
        if err != nil {
            log.Panic("Parsing of float went bad: ", err)
        }
    }

    return output

}



