package gsea

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type GeneSet struct {
	Name  string
	Url   string
	Genes []string
}

// Parses a GMT file from msigdb
// (http://software.broadinstitute.org/gsea/downloads.jsp) and returns the gene
// sets found.
func ParseGMT(filename string) ([]GeneSet, error) {

	file, err := os.Open(filename)
	if err != nil {
		return []GeneSet{}, errors.Wrap(err, "Could not open GMT file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	geneSets := []GeneSet{}

	for scanner.Scan() {
		line := scanner.Text()
		columns := strings.SplitN(line, "\t", 3)
		if len(columns) < 3 {
			return []GeneSet{}, errors.New("Could not parse GMT file")
		}
		gs := GeneSet{
			Name:  columns[0],
			Url:   columns[1],
			Genes: strings.Split(columns[2], "\t")}
		geneSets = append(geneSets, gs)
	}

	if err := scanner.Err(); err != nil {
		return []GeneSet{}, errors.Wrap(err, "Error reading GMT file")
	}

	return geneSets, nil
}

// Write genesets to GMT file following the MSigDB format.
func WriteGMT(genesets []GeneSet, filename string) error {

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	for _, geneset := range genesets {
		genes := strings.Join(geneset.Genes, "\t")
		line := strings.Join([]string{geneset.Name, geneset.Url, genes}, "\t")
		fmt.Fprintln(w, line)
	}

	return w.Flush()
}
