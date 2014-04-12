package main

import (
    "testing" 
    "strconv"
    "log"
    "io/ioutil"
) 

var ds Dataset
var restService *RestService

func BenchmarkInit(b *testing.B) { 
    disableLogging()
    path := "/Users/bjorn/stallo/src/src/nowac/datastore/data"
    // Set up rest service to be used later. reset timer afterwards
    restService = Init(path)
    b.ResetTimer()

    for n := 0; n < b.N; n++ {
        ds = NewDataset(path) 
    } 
} 


var result []float64
func BenchmarkGeneExpression(b *testing.B) {
    for n := 0; n < b.N; n++ {
        result = restService.GeneExpression("hsa:23533")
    }
}


var res float64
func BenchmarkAverageDiff(b *testing.B) {
    for n := 0; n < b.N; n++ {
        res = restService.AvgDiff("hsa:23533")
    }
}

func BenchmarkStd(b *testing.B) {
    for n := 0; n < b.N; n++ {
        res = restService.Std("hsa:23533")
    }
}

func BenchmarkVariance(b *testing.B) {
    for n := 0; n < b.N; n++ {
        res = restService.Variance("hsa:23533")
    }
}

func BenchmarkSetScale(b *testing.B){ 
    for n := 0; n < b.N; n++ {
        restService.SetScale("log")
    } 
} 

var exprs, bg string
func BenchmarkBackground(b *testing.B) {
    // get expression value from prev test 
    exprs = strconv.FormatFloat(result[0],'f',-1, 64)

    for n := 0; n < b.N; n++ {
        bg = restService.Bg("hsa:23533", exprs)
    }
}

func disableLogging() {
    log.SetOutput(ioutil.Discard) 
} 


