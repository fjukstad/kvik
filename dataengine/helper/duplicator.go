package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	. "os"
	"strconv"
	"time"
)

func randomFloat() float64 {
	// seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomFloat := r.Float64()
	return randomFloat
}

func randomInt(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(n)
	res := strconv.Itoa(num)
	return res
}

func randomFloats(num int) []string {
	var res []string
	for i := 0; i < num; i++ {
		f := strconv.FormatFloat(randomFloat()*1000, 'f', -1, 64)
		res = append(res, f)
	}
	return res
}

func main() {
	var numPairs = flag.Int("numpairs", 77, "number of cc pairs")

	var exprsFilename = flag.String("exprsFilename", "exprs.csv",
		"output filename for expression data ")
	var bgFilename = flag.String("bgFilename", "background.csv",
		"output filename for background data")

	var dsfile = flag.String("dsfile", "/Users/bjorn/stallo/data/exprs.csv",
		"original dataset. where to get gene names")
	flag.Parse()

	// Read header from original dataset
	exprsfile, err := os.Open(*dsfile)

	if err != nil {
		log.Panic(err)
	}

	defer exprsfile.Close()

	reader := csv.NewReader(exprsfile)
	header, err := reader.Read()

	numGenes := len(header) - 1

	exprsFile, err := os.OpenFile(*exprsFilename, O_RDWR|O_CREATE, 0666)

	if err != nil {
		log.Println("Could not open or create file", err)
		return
	}

	bgFile, err := os.OpenFile(*bgFilename, O_RDWR|O_CREATE, 0666)

	if err != nil {
		log.Println("Could not open or create file", err)
		return
	}

	log.Println("Writing dataset with ", numGenes, "genes and", *numPairs,
		"cc pairs to", *exprsFilename, "and background information to",
		*bgFilename)

	exprsWriter := csv.NewWriter(exprsFile)
	bgWriter := csv.NewWriter(bgFile)

	// write header to file
	exprsWriter.Write(header)

	//write bg header to bg file
	bgHeader := []string{"labnr", "id", "Case_ctrl", "ET", "EN", "stage"}
	bgWriter.Write(bgHeader)

	var np int
	np = *numPairs
	for i := 0; i < *numPairs; i++ {
		p := strconv.Itoa(i)

		record := make([]string, 0)
		record = append(record, p)
		record = append(record, randomFloats(numGenes)...)

		exprsWriter.Write(record)
		if err := exprsWriter.Error(); err != nil {
			log.Println(err)
		}

		labnr := randomInt(200)
		id := p
		cc := "\"case\""
		// these numbers below are randomly picked
		et := randomInt(5)
		en := randomInt(5)
		stage := randomInt(3)

		bgrec := []string{labnr, id, cc, et, en, stage}
		bgWriter.Write(bgrec)

		p = p + "_1"
		record = make([]string, 0)
		record = append(record, p)
		record = append(record, randomFloats(numGenes)...)
		exprsWriter.Write(record)

		if err := exprsWriter.Error(); err != nil {
			log.Println(err)
		}

		labnr = randomInt(200)
		id = p
		cc = "ctrl"
		// NA for ctrl
		et = "NA"
		en = "NA"
		stage = "NA"

		bgrec = []string{labnr, id, cc, et, en, stage}
		bgWriter.Write(bgrec)

		fmt.Print(i, " of ", np, " iterations done \r")
	}
	fmt.Print("\n")

	exprsWriter.Flush()
	bgWriter.Flush()

	if err := exprsWriter.Error(); err != nil {
		log.Println(err)
	}
	if err := bgWriter.Error(); err != nil {
		log.Println(err)
	}
}

// id, gene, gene, gene
// integer, gene, gene, gene
// integer_1, gene, gene, gene
// integer2, gene, gene, gene
// integer2_1, gene, gene, gene
