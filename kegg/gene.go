package kegg

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/fjukstad/gocache"
)

type Gene struct {
	Id          string
	Name        string
	Definition  string
	Orthology   string
	Organism    string
	Pathways    []string
	Modules     []string
	Diseases    []string
	Drug_Target []string
	Classes     []string
	Position    string
	Motif       string
	DBLinks     map[string]string
	Structure   string
	AASEQ       Sequence
	NTSEQ       Sequence
}

type Sequence struct {
	length   int
	Sequence string
}

func GetAllHumanGenes() []Gene {

	url := "http://rest.kegg.jp/list/hsa"
	resp, err := gocache.Get(url)
	if err != nil {
		log.Panic("Could not fetch gene list", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic("Could not read body of gene list response", err)
	}

	var geneids []string

	lines := strings.Split(string(body), "\n")

	for _, line := range lines {
		id := strings.Split(line, "\t")[0]
		geneid := strings.Split(id, ":")

		if len(geneid) < 2 {
			continue
		}

		geneids = append(geneids, geneid[1])
	}

	var result []Gene
	for _, id := range geneids {
		result = append(result, GetGene(id))
	}

	return result

}

func GetGene(id string) Gene {
	baseURL := "http://rest.kegg.jp/get/hsa:"
	url := baseURL + id
	log.Println(url)
	response, err := gocache.Get(url)
	if err != nil {
		log.Panic("Cannot download from url:", err)
	}

	/*
	   body, err := ioutil.ReadAll(response.Body)
	   if err != nil{
	       log.Panic("Error reading body:",err)
	   }
	*/

	gene := parseGeneResponse(response.Body)

	return gene
}

type GenePathways struct {
	GeneId   string
	Pathways []string
}

// Fetch pathways for a specific gene
func Pathways(g Gene) GenePathways {
	pws := GenePathways{g.Id, g.Pathways}
	return pws

}

// TODO:

func PathwaysJSON(pws GenePathways) string {
	b, err := json.Marshal(pws)
	if err != nil {
		log.Panic("marshaling went bad: ", err)
	}
	return string(b)
}

func GeneJSON(g Gene) string {

	b, err := json.Marshal(g)
	if err != nil {
		log.Panic("Marshaling went bad: ", err)
	}

	return string(b)

}

func (g Gene) Print() {
	fmt.Println("\nGene")
	fmt.Println("\tId:", g.Id)
	fmt.Println("\tName:", g.Name)
	fmt.Println("\tDefinition:", g.Definition)
	fmt.Println("\tOrthology:", g.Orthology)
	fmt.Println("\tOrganism:", g.Organism)
	fmt.Println("\tPathways:", g.Pathways)
	fmt.Println("\tClasses:", g.Classes)
	fmt.Println("\tPosition:", g.Position)
	fmt.Println("\tMotif:", g.Motif)
	fmt.Println("\tDBLinks:", g.DBLinks)
	fmt.Println("\tStructure:", g.Structure)
	//fmt.Println("\tAASEQ:", g.AASEQ)
	//fmt.Println("\tNTSEQ:", g.NTSEQ)
	fmt.Println("")
}

func GeneIdFromName(name string) string {
	baseURL := "http://rest.kegg.jp/find/genes/"
	url := baseURL + name

	response, err := gocache.Get(url)
	if err != nil {
		log.Panic(err)
	}

	tsv := csv.NewReader(response.Body)
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

	gene := strings.Split(records[0][0], " ")[0]

	//id := strings.TrimLeft(gene, "hsa:")

	return gene

}

func parseGeneResponse(response io.ReadCloser) Gene {

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
	gene := Gene{}

	tmp := make([]string, 0)
	current := ""
	sequence := ""
	for i := range records {

		line := strings.Split(records[i][0], " ")

		switch line[0] {
		case "ENTRY":
			// parsing to extract id
			a := strings.Join(line[7:], " ")
			b := strings.Split(a, " ")[0]
			gene.Id = b

		case "NAME":
			gene.Name = strings.Join(line[8:], " ")

		case "DEFINITION":
			gene.Definition = strings.Join(line[2:], " ")

		case "ORTHOLOGY":
			gene.Orthology = strings.Join(line[3:], " ")

		case "ORGANISM":
			gene.Organism = strings.Join(line[4:], " ")

		case "PATHWAY":
			current = "PATHWAY"
			// Parsing to extract the hsa12345 string
			a := strings.Join(line[5:], " ")
			b := strings.Split(a, " ")[0]
			tmp = append(tmp, b)

		// NOTE: We're currently not storing the BRITE field
		case "BRITE":

			if current == "PATHWAY" {
				gene.Pathways = tmp
			}

			current = "BRITE"

		case "DISEASE":

			if current == "PATHWAY" {
				gene.Pathways = tmp
			}
			if current == "MODULE" {
				gene.Modules = tmp
			}
			current = "DISEASE"
			tmp = make([]string, 0)
			a := strings.Join(line[5:], " ")
			tmp = append(tmp, a)

		case "MODULE":
			if current == "PATHWAY" {
				gene.Pathways = tmp
			}
			if current == "DISEASE" {
				gene.Diseases = tmp
			}

			current = "MODULE"
			tmp = make([]string, 0)
			a := strings.Join(line[6:], " ")
			tmp = append(tmp, a)

		case "DRUG_TARGET":
			if current == "DISEASE" {
				gene.Diseases = tmp
			}
			current = "DRUG_TARGET"
			tmp = make([]string, 0)
			a := strings.Join(line[1:], " ")
			tmp = append(tmp, a)

		case "CLASS":
			if current == "PATHWAY" {
				gene.Pathways = tmp
			}
			if current == "DISEASE" {
				gene.Diseases = tmp
			}

			if current == "MODULE" {
				gene.Modules = tmp
			}
			if current == "DRUG_TARGET" {
				gene.Drug_Target = tmp
			}

			current = "CLASS"
			tmp = make([]string, 0)
			a := strings.Join(line[7:], " ")
			tmp = append(tmp, a)

		case "POSITION":
			current = "POSITION"
			gene.Classes = tmp
			gene.Position = strings.Join(line[4:], " ")

		case "MOTIF":
			gene.Motif = strings.Join(line[7:], " ")

		case "DBLINKS":
			current = "DBLINKS"
			tmp = make([]string, 0)
			a := strings.Join(line[5:], " ")
			tmp = append(tmp, a)
			e := strings.Split(a, ":")
			gene.DBLinks = make(map[string]string)
			link := strings.Replace(e[1], " ", "", -1)
			gene.DBLinks[e[0]] = link

		case "STRUCTURE":
			current = "STRUCTURE"
			gene.Structure = strings.Join(line[3:], " ")

		case "AASEQ":
			current = "AASEQ"
			a := strings.Join(line[7:], " ")
			length, err := strconv.Atoi(a)
			if err != nil {
				log.Panic("AASEQ PARSING ERROR:", err)
			}
			gene.AASEQ = Sequence{length, ""}

		case "NTSEQ":
			gene.AASEQ.Sequence = sequence
			current = "NTSEQ"
			sequence = ""
			a := strings.Join(line[7:], " ")
			length, err := strconv.Atoi(a)
			if err != nil {
				log.Panic("NTSEQ PARSING ERROR:", err)
			}
			gene.NTSEQ = Sequence{length, ""}

		case "///":
			gene.NTSEQ.Sequence = sequence

		default:
			if current == "PATHWAY" {
				// Again some fancy parsing to extract hsa1234 string
				a := strings.Join(line[12:], " ")
				b := strings.Split(a, " ")[0]
				tmp = append(tmp, b)
			}
			if current == "CLASS" ||
				current == "MODULE" ||
				current == "DRUG_TARGET" {

				// Parsing, not very pretty...
				a := strings.Join(line[0:], " ")
				b := strings.Replace(a, "    ", "", -1)
				tmp = append(tmp, b)
			}
			if current == "DBLINKS" {
				a := strings.Join(line[0:], " ")
				b := strings.Replace(a, "    ", "", -1)
				e := strings.Split(b, ":")

				link := strings.Replace(e[1], " ", "", -1)
				gene.DBLinks[e[0]] = link
			}
			if current == "AASEQ" ||
				current == "NTSEQ" {
				a := strings.Join(line[0:], " ")
				b := strings.Replace(a, "    ", "", -1)
				sequence += b
			}
		}
	}

	return gene
}
