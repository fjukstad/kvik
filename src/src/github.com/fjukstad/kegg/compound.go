package kegg

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/fjukstad/gocache"
)

type Compound struct {
	Entry      string
	Name       []string
	Formula    string
	Exact_Mass string
	Mol_Weight string
	Remark     string
	Comment    string
	Reaction   []string
	Pathway    []string
	Module     string
	Enzyme     []string
	Brite      string
	DBLinks    map[string]string
	Atom       []string
	Bond       []string
}

func GetCompound(id string) Compound {
	baseURL := "http://rest.kegg.jp/get/cpd:"
	url := baseURL + id

	response, err := gocache.Get(url)
	if err != nil {
		log.Panic("Cannot download from url:", err)
	}

	compound := parseCompoundResponse(response.Body)

	return compound

}

func (c Compound) Print() {
	fmt.Println("\nCompound")
	fmt.Println("\tEntry:", c.Entry)
	fmt.Println("\tName:", c.Name)
	fmt.Println("\tFormula:", c.Formula)
	fmt.Println("\tExact mass:", c.Exact_Mass)
	fmt.Println("\tMol weight:", c.Mol_Weight)
	fmt.Println("\tRemark:", c.Remark)
	fmt.Println("\tComment:", c.Comment)
	fmt.Println("\tReaction:", c.Reaction)
	fmt.Println("\tPathway:", c.Pathway)
	fmt.Println("\tModule:", c.Module)
	fmt.Println("\tEnzyme:", c.Enzyme)
	fmt.Println("\tBrite:", c.Brite)
	fmt.Println("\tDBLinks:", c.DBLinks)
	fmt.Println("\tAtom:", c.Atom)
	fmt.Println("\tBond:", c.Bond)
	return
}

func (c Compound) JSON() string {
	b, err := json.Marshal(c)
	if err != nil {
		log.Panic("Marshaling went horrible :( ", err)
	}
	return string(b)
}

func parseCompoundResponse(response io.ReadCloser) Compound {
	tsv := csv.NewReader(response)
	tsv.Comma = '\t'
	tsv.Comment = '#'
	tsv.LazyQuotes = true
	tsv.TrailingComma = true
	tsv.TrimLeadingSpace = false
	tsv.FieldsPerRecord = -1

	records, err := tsv.ReadAll()

	if err != nil {
		log.Println(response)
		log.Panic("Error reading records:", err, records)

	}

	compound := Compound{}
	current := ""

	// Nice figure to put in vis:
	// http://www.genome.jp/Fig/compound/compound.Entry.gif

	for _, rec := range records {

		line := strings.Split(rec[0], " ")
		switch line[0] {
		case "ENTRY":
			a := strings.Join(line[7:], " ")
			b := strings.Split(a, " ")[0]
			compound.Entry = b
		case "NAME":
			a := strings.Join(line[8:], " ")
			b := strings.Replace(a, ";", "", -1)
			compound.Name = append(compound.Name, b)
			current = "NAME"
		case "FORMULA":
			current = "FORMULA"
			a := strings.Join(line[5:], " ")
			compound.Formula = a
		case "EXACT_MASS":
			a := strings.Join(line[2:], " ")
			compound.Exact_Mass = a
		case "MOL_WEIGHT":
			a := strings.Join(line[2:], " ")
			compound.Mol_Weight = a
		case "REACTION":
			a := strings.Join(line[2:], " ")
			b := strings.Split(a, " ")
			for i := 0; i < len(b); i++ {
				if i > 1 {
					compound.Reaction = append(compound.Reaction, b[i])
				}
			}
		case "PATHWAY":
			a := strings.Join(line[5:], " ")
			b := strings.Split(a, " ")[0]
			compound.Pathway = append(compound.Pathway, b)
			current = "PATHWAY"
		case "ENZYME":
			current = "ENZYME"
		case "BRITE":
			current = "BRITE"
		case "DBLINKS":
			a := strings.Join(line[5:], " ")
			b := strings.Split(a, ":")
			compound.DBLinks = make(map[string]string)
			compound.DBLinks[b[0]] = b[1]
			current = "DBLINKS"
		case "ATOM":
			current = "ATOM"
		case "BOND":
			current = "BOND"
		case "MODULE":
			current = "MODULE"
		case "  ORGANISM":
			current = "ORGANISM"

		default:
			if current == "NAME" {
				a := strings.Join(line[0:], " ")
				b := strings.Replace(a, "    ", "", -1)
				c := strings.Replace(b, ";", "", -1)
				compound.Name = append(compound.Name, c)
			}

			if current == "PATHWAY" {
				a := strings.Join(line[12:], " ")
				b := strings.Split(a, " ")[0]
				compound.Pathway = append(compound.Pathway, b)
			}

			if current == "DBLINKS" {
				// if line is not "///" last line
				if len(line) > 5 {
					a := strings.Join(line[12:], " ")
					b := strings.Split(a, ":")
					compound.DBLinks[b[0]] = b[1]
				}

			}
		}

	}

	compound.Print()
	return compound

}
