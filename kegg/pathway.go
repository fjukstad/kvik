package kegg

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/fjukstad/gocache"
)

type Pathway struct {
	Id          string
	Name        string
	Description string
	Class       string
	Pathway_Map string
	Diseases    []string
	Drugs       []string
	DBLinks     []string
	Organism    string
	Genes       []string
	Compounds   []string
}

type KeggPathway struct {
	XMLName   xml.Name       `xml:"pathway"`
	Name      string         `xml:"name,attr"`
	Org       string         `xml:"org,attr"`
	Number    string         `xml:"number,attr"`
	Title     string         `xml:"title,attr"`
	Image     string         `xml:"image,attr"`
	Link      string         `xml:"link,attr"`
	Entries   []KeggEntry    `xml:"entry"`
	Relations []KeggRelation `xml:"relation"`
}

type KeggEntry struct {
	Id       string       `xml:"id,attr"`
	Name     string       `xml:"name,attr"`
	Type     string       `xml:"type,attr"`
	Link     string       `xml:"link,attr"`
	Graphics KeggGraphics `xml:"graphics"`
}

type KeggRelation struct {
	Entry1   string        `xml:"entry1,attr"`
	Entry2   string        `xml:"entry2,attr"`
	Type     string        `xml:"type,attr"`
	Subtypes []KeggSubtype `xml:"relation>subtype"`
}

type KeggGraphics struct {
	Name    string `xml:"name,attr"`
	Fgcolor string `xml:"fgcolor,attr"`
	Bgcolor string `xml:"bgcolor,attr"`
	Type    string `xml:"type,attr"`
	X       string `xml:"x,attr"`
	Y       string `xml:"y,attr"`
	Width   string `xml:"width,attr"`
	Height  string `xml:"height,attr"`
}

type KeggSubtype struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func GetAllHumanPathways() []Pathway {

	url := "http://rest.kegg.jp/list/pathway/hsa"
	resp, err := gocache.Get(url)
	if err != nil {
		log.Panic("Could not fetch pathway list", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic("Could not read body of pathway list response", err)
	}

	var res []string

	lines := strings.Split(string(body), "\n")
	for i := range lines {
		line := lines[i]
		id := strings.Split(line, "\t")[0]
		pathid := strings.Split(id, ":")

		// empty line
		if len(pathid) < 2 {
			continue
		}

		pwid := pathid[1]
		res = append(res, pwid)

	}

	// sort list of pathways
	res = SortPathwayIds(res)

	// Get all info
	var pathways []Pathway
	//:= make([]Pathway, 0)
	for _, id := range res {
		pathways = append(pathways, GetPathway(id))
	}

	return pathways

}

type ByName []Pathway

func (a ByName) Len() int {
	return len(a)
}
func (a ByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

func SortPathwayIds(ids []string) []string {
	pws := make([]Pathway, 0)
	for _, id := range ids {
		pws = append(pws, GetPathway(id))
	}

	sort.Sort(ByName(pws))

	pwids := make([]string, 0)
	for _, pw := range pws {
		pwids = append(pwids, pw.Id)
	}

	return pwids

}

func ReadablePathwayNames(ids []string) []string {

	pathways := make([]string, len(ids))

	for i, id := range ids {
		pathways[i] = ReadablePathwayName(id)
	}

	return pathways

}

func ReadablePathwayName(id string) string {
	//
	name := GetPathway(id).Name
	shortName := strings.Trim(strings.SplitAfter(name, " - ")[0], " - ")

	return shortName
}

func getMap(url string) []byte {
	response, err := gocache.Get(url)
	if err != nil {
		log.Panic("Could not download pathway kgml:", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Panic("KGML body could not be read:", err)
	}
	return body
}

func (pathway *KeggPathway) Print() {
	log.Println("XMLName:", pathway.XMLName)
	log.Println("Name:", pathway.Name)
	log.Println("Org:", pathway.Org)
	log.Println("Title:", pathway.Title)
	log.Println("Image:", pathway.Image)
	log.Println("Link", pathway.Link)
	log.Println("Entries:", pathway.Entries)
	log.Println("Relations:", pathway.Relations)
}

func NewKeggPathway(keggId string) *KeggPathway {

	baseURL := "http://rest.kegg.jp/get/" + keggId
	url := baseURL + "/kgml"

	// url := "http://localhost:8000/public/pathway.kgml"
	pw := getMap(url)

	pathway := new(KeggPathway)
	err := xml.Unmarshal(pw, pathway)

	if err != nil {
		log.Panic("Could not unmarshal KGML ", err)
	}

	return pathway

}

type Node struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Type            int    `json:"type"`
	Size            int    `json:"size"`
	Description     string `json:"description"`
	ForegroundColor string `json:"fgcolor"`
	BackgroundColor string `json:"bgcolor"`
	Shape           string `json:"shape"`
	X               int    `json:"x"`
	Y               int    `json:"y"`
	Height          int    `json:"height"`
	Width           int    `json:"width"`
}

type Edge struct {
	Source int `json:"source"`
	Target int `json:"target"`
	Index  int `json:"index"`
	Weight int `json:"weight"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

func PathwayGraph(keggId string) Graph {

	baseURL := "http://rest.kegg.jp/get/" + keggId
	url := baseURL + "/kgml"

	pw := getMap(url)

	pathway := new(KeggPathway)
	err := xml.Unmarshal(pw, pathway)

	if err != nil {
		log.Print("Could not unmarshal KGML ", keggId, err)
		pw := getMap(url)
		err = xml.Unmarshal(pw, pathway)
		if err != nil {
			log.Panic("Could not unmarshal KGML for the second time", keggId, err)
		}
	}

	imgurl := "http://www.genome.jp/kegg/pathway/hsa/" + keggId + ".png"

	resp, err := http.Get(imgurl)
	if err != nil {
		log.Panic("Image could not be downloaded ", err)
	}

	img, err := png.Decode(resp.Body)

	if err != nil {
		log.Panic("Image could not be decoded ", err)
	}

	imgrect := img.Bounds()

	sizeX := imgrect.Max.X - imgrect.Min.X
	sizeY := imgrect.Max.Y - imgrect.Min.Y

	// Store image for later use
	path := "public/pathways/"
	filename := keggId + ".png"
	err = storeImage(path, filename, img)
	if err != nil {
		log.Panic("Image could not be stored", err)
	}

	// First create a node that will serve as a background image to the pathway
	var background = Node{
		0,
		"bg",
		0,
		1,
		imgurl,
		"#fff",
		"#fff",
		"rectangle",
		sizeX / 2,
		sizeY / 2,
		sizeY,
		sizeX}

	var node Node
	var edge Edge
	var nodes []Node
	var edges []Edge

	nodes = append(nodes, background)

	// Generate some nodes
	for j := range pathway.Entries {
		ent := pathway.Entries[j]
		id, _ := strconv.Atoi(ent.Id)
		name := ent.Name
		t, _ := strconv.Atoi(ent.Type)
		size := 1
		graphics := ent.Graphics

		// Trimming away :title for the node containing the pathway name
		description := strings.TrimPrefix(strings.Split(graphics.Name, ",")[0], "TITLE:")
		fgcolor := graphics.Fgcolor
		bgcolor := graphics.Bgcolor
		shape := graphics.Type
		x, _ := strconv.Atoi(graphics.X)
		y, _ := strconv.Atoi(graphics.Y)
		height, _ := strconv.Atoi(graphics.Height)
		width, _ := strconv.Atoi(graphics.Width)

		node = Node{id, name, t, size, description, fgcolor, bgcolor, shape,
			x, y, height, width}

		nodes = append(nodes, node)
	}

	// Generate some edges
	for i := range pathway.Relations {
		rel := pathway.Relations[i]
		source, _ := strconv.Atoi(rel.Entry1)
		target, _ := strconv.Atoi(rel.Entry2)
		weight := 19
		edge = Edge{source, target, i, weight}
		edges = append(edges, edge)
	}

	graph := Graph{nodes, edges}

	return graph

}

func storeImage(path, filename string, image image.Image) error {

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	fn := path + "/" + filename
	file, err := os.Create(fn)
	if err != nil {
		return err
	}

	return png.Encode(file, image)

}

func GetPathway(id string) Pathway {
	baseURL := "http://rest.kegg.jp/get/"
	url := baseURL + id

	response, err := gocache.Get(url)
	if err != nil {
		log.Panic("Cannot download pathway:", err)
	}

	pathway := parsePathwayResponse(response.Body)

	return pathway

}

func (p Pathway) JSON() string {
	b, err := json.Marshal(p)
	if err != nil {
		log.Panic("marshaling went bad: ", err)
	}
	return string(b)

}

func (p Pathway) Print() {
	fmt.Println("\nPathway")
	fmt.Println("\tId:", p.Id)
	fmt.Println("\tName:", p.Name)
	fmt.Println("\tClass:", p.Class)
	fmt.Println("\tPathway map:", p.Pathway_Map)
	fmt.Println("\tDiseases:", p.Diseases)
	fmt.Println("\tDrugs:", p.Drugs)
	fmt.Println("\tOrganism:", p.Organism)
	fmt.Println("\tGenes:", p.Genes)
	fmt.Println("\tCompounds:", p.Compounds)
	fmt.Println("")
}

func GiveMeSomePathways() []string {
	pw := []string{
		"hsa05200",
		"hsa04915",
		"hsa04612",
		"hsa04062",
		"hsa04660",
		"hsa04630",
		"hsa04151",
		"hsa04310",
		"hsa04662",
	}
	return pw
}
func parsePathwayResponse(response io.ReadCloser) Pathway {

	tsv := csv.NewReader(response)
	tsv.Comma = '\t'
	tsv.Comment = '#'
	tsv.LazyQuotes = true
	tsv.TrailingComma = true
	tsv.TrimLeadingSpace = false

	records, err := tsv.ReadAll()

	if err != nil {
		log.Panic("Error reading records:", err)
	}

	p := Pathway{}
	tmp := make([]string, 0)
	tmpstring := ""
	current := ""
	for i := range records {

		line := strings.Split(records[i][0], " ")

		switch line[0] {
		case "ENTRY":
			a := strings.Join(line[7:], " ")
			b := strings.Split(a, " ")[0]
			p.Id = b

		case "NAME":
			p.Name = strings.Join(line[8:], " ")

		case "DESCRIPTION":
			current = "DESCRIPTION"
			tmpstring = strings.Join(line[1:], " ")

		case "CLASS":
			if current == "DESCRIPTION" {
				p.Description = tmpstring
			}

			current = "CLASS"
			p.Class = strings.Join(line[8:], " ")

		case "PATHWAY_MAP":
			p.Pathway_Map = strings.Join(line[1:], " ")

		case "DISEASE":
			current = "DISEASE"

			a := strings.Join(line[5:], " ")
			tmp = append(tmp, a)

		case "DBLINKS":
			p.Diseases = tmp
			current = "DBLINKS"

			a := strings.Join(line[5:], " ")
			tmp = make([]string, 0)
			tmp = append(tmp, a)

		case "DRUG":
			p.Diseases = tmp
			tmp = make([]string, 0)
			current = "DRUG"
			a := strings.Join(line[8:], " ")
			tmp = append(tmp, a)

		case "ORGANISM":
			if current == "DISEASE" {
				p.Drugs = tmp
			}
			if current == "DBLINKS" {
				p.DBLinks = tmp
			}
			current = "ORGANISM"
			p.Organism = strings.Join(line[4:], " ")

		case "GENE":
			current = "GENE"
			tmp = make([]string, 0)
			a := strings.Join(line[8:], " ")
			b := strings.Split(a, " ")[0]
			tmp = append(tmp, b)

		case "COMPOUND":
			p.Genes = tmp
			current = "COMPOUND"
			tmp = make([]string, 0)
			a := strings.Join(line[4:], " ")

			tmp = append(tmp, a)

		case "REFERENCE":
			p.Compounds = tmp
			current = "REFERENCE"
			break

		default:
			if current == "DESCRIPTION" {
				tmpstring = tmpstring + strings.Join(line[1:], " ")
			}

			if current == "DISEASE" {
				a := strings.Join(line[12:], " ")
				tmp = append(tmp, a)
			}
			if current == "DRUG" {
				a := strings.Join(line[0:], " ")
				b := strings.Replace(a, "    ", "", -1)
				tmp = append(tmp, b)
			}
			if current == "GENE" {
				a := strings.Join(line[0:], " ")
				a = strings.Replace(a, "    ", "", -1)
				b := strings.Split(a, " ")[0]
				tmp = append(tmp, b)

			}
			if current == "COMPOUND" {
				a := strings.Join(line[0:], " ")
				a = strings.Replace(a, "    ", "", -1)

				tmp = append(tmp, a)

			}

		}
	}
	return p

}
