package gsea

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type GeneSet struct {
	Name  string
	Url   string
	Genes []string
}

func ParseGMT(filename string) ([]GeneSet, error) {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
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
