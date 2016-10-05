package main

import "github.com/fjukstad/kvik/gsea"
import "fmt"

func main() {
	abstract, err := gsea.Abstract("BIOCARTA_HER2_PATHWAY")
	fmt.Println(abstract, err)

	brief, err := gsea.BriefDescription("PID_HES_HEY_PATHWAY")
	fmt.Println(brief, err)

	systematicName, err := gsea.SystematicName("PID_HES_HEY_PATHWAY")
	fmt.Println(systematicName, err)

	collection, err := gsea.Collection("PID_HES_HEY_PATHWAY")
	fmt.Println(collection, len(collection), err)

	pub, err := gsea.SourcePublication("PID_HES_HEY_PATHWAY")
	fmt.Println(pub, err)

	fmt.Println(gsea.PublicationURL(pub))

	related, err := gsea.RelatedGeneSets("PID_HES_HEY_PATHWAY")
	fmt.Println(related, err)

	org, err := gsea.Organism("PID_HES_HEY_PATHWAY")
	fmt.Println(org, err)

	contrib, err := gsea.ContributedBy("PID_HES_HEY_PATHWAY")
	fmt.Println(contrib, err)

	platform, err := gsea.SourcePlatform("PID_HES_HEY_PATHWAY")
	fmt.Println(platform, err)

	cellLines := gsea.CompendiumURL("PID_HES_HEY_PATHWAY", "cancerCellLines")
	fmt.Println(cellLines)

	novartisHuman := gsea.CompendiumURL("PID_HES_HEY_PATHWAY", "novartisHuman")
	fmt.Println(novartisHuman)

	dsref, err := gsea.DatasetReference("GSE3982_CTRL_VS_LPS_4H_MAC_UP")
	fmt.Println(dsref, err)

}
