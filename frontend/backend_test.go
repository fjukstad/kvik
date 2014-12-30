package main 

import ( 
    "testing"
    "log" 
    "io/ioutil"
    "os/exec" 
) 

var rs *NOWACService

var res string 

func BenchmarkNewPathwayGraph(b *testing.B) { 
    disableLogging()
    rs = new(NOWACService) 
    res = rs.NewPathwayGraph("pathwaySelect=hsa04915") 
    
    b.ResetTimer()

    for n := 0; n < b.N; n++ {
        res = rs.NewPathwayGraph("pathwaySelect=hsa04915") 
    } 
} 

func BenchmarkNewPathwayGraphNoCache(b *testing.B) { 
    clearCache()
    b.ResetTimer()
    for n := 0; n < b.N; n++ {
        res = rs.NewPathwayGraph("pathwaySelect=hsa04915") 
        Clear(b) 
    } 
} 


func BenchmarkGetInfo(b *testing.B) { 
    res = rs.GetInfo("hsa:2776", "unused arg")
    b.ResetTimer()

    for n := 0; n < b.N; n++ {
        res = rs.GetInfo("hsa:2776", "unused arg")
    } 
} 

func BenchmarkGetInfoNoCache(b *testing.B) { 
    Clear(b) 
    for n := 0; n < b.N; n++ {
        res = rs.GetInfo("hsa:2776", "unused arg")
        Clear(b)
    } 
} 


func BenchmarkGetGeneVis(b *testing.B) { 
    res = rs.GetGeneVis("2776")
    b.ResetTimer()

    for n := 0; n < b.N; n++ {
        res = rs.GetGeneVis("2776")
    } 
} 

/*
func BenchmarkGetGeneVisNoCache(b *testing.B) { 
    Clear(b)
    for n := 0; n < b.N; n++ {
        res = rs.GetGeneVis("2776")
        Clear(b)
    } 
} 
*/

func BenchmarkPathways(b *testing.B) { 
    res = rs.Pathways("hsa:2776")
    b.ResetTimer()

    for n := 0; n < b.N; n++ {
        res = rs.Pathways("hsa:2776")
    } 
} 

func BenchmarkPathwaysNoCache(b *testing.B) { 
    Clear(b)
    for n := 0; n < b.N; n++ {
        res = rs.Pathways("hsa:2776")
        Clear(b)
    } 
} 

/*
func BenchmarkCommonPathways(b *testing.B) { 
    res = rs.CommonPathways("hsa:2776")
    b.ResetTimer()

    for n := 0; n < b.N; n++ {
        res = rs.CommonPathways("hsa:2776")
    } 
} 

func BenchmarkCommonPathwaysNoCache(b *testing.B) { 
    Clear(b)
    for n := 0; n < b.N; n++ {
        res = rs.CommonPathways("hsa:2776")
        Clear(b)
    } 
} 
*/




func BenchmarkPathwayGeneCount(b *testing.B) { 
    res = rs.PathwayGeneCount("hsa:2776")
    b.ResetTimer()

    for n := 0; n < b.N; n++ {
        res = rs.PathwayGeneCount("hsa:2776")
    } 
} 


func BenchmarkPathwayGeneCountNoCache(b *testing.B) { 
    Clear(b)
    for n := 0; n < b.N; n++ {
        res = rs.PathwayGeneCount("hsa:2776")
        Clear(b)
    } 
} 


var a int
func BenchmarkCommonGenes(b *testing.B) { 
    a = rs.CommonGenes("hsa04915+hsa04921")
    b.ResetTimer()

    for n := 0; n < b.N; n++ {
        a = rs.CommonGenes("hsa04915+hsa04921")
    } 
} 

func BenchmarkCommonGenesNoCache(b *testing.B) { 
    Clear(b)
    for n := 0; n < b.N; n++ {
        a = rs.CommonGenes("hsa04915+hsa04921")
        Clear(b)
    } 
} 





func disableLogging() {
    log.SetOutput(ioutil.Discard) 
} 

func Clear(b *testing.B) {
        b.StopTimer()
        clearCache()
        b.StartTimer()
} 


func clearCache() { 
    cmd := exec.Command("rm", "-rf", "cache") 
    err := cmd.Run() 
    log.Println(err) 
} 
