package eutils

type SummaryResult struct {
	DocumentSummarySet DocumentSummarySet
}

type DocumentSummarySet struct {
	DbBuild         string
	Status          string `xml:"status,attr"`
	DocumentSummary []DocumentSummary
}

type DocumentSummary struct {
	Uid             string `xml:"uid,attr"`
	PubDate         string
	Source          string
	Authors         []Authors
	Title           string
	SortTitle       string
	Volume          string
	Issue           string
	Pages           string
	Lang            string `xml:"Lang>string"`
	NlmUniqueID     int
	ISSN            string
	ESSN            string
	PubType         string `xml:"PubType>flag"`
	RecordStatus    string
	PubStatus       int
	ArticleIds      []ArticleId
	History         []PubMedPubDate
	Attributes      string `xml:"Attributes>flag"`
	PmcRefCount     int
	FullJournalName string
	ViewCount       int
	DocType         string
	SortPubDate     string
	SortFirstAuthor string

	Name              string
	Description       string
	Status            int
	CurrentId         int
	Chromosome        int
	GeneticSource     string
	OtherAliases      string
	OtherDesignations string
	Mim               int `xml:"Mim>int"`
	GenomicInfo
	Summary  string
	ChrSort  int
	ChrStart int
	Organism
	LocationHist []LocationHistType
}
