package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/fjukstad/kvik/eutils"
	"github.com/fjukstad/kvik/genenames"
)

func DocHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	if len(values["geneSymbol"]) < 1 {
		http.Error(w, "No gene symbol specified", http.StatusBadRequest)
		return
	}
	geneSymbol := values["geneSymbol"][0]

	// cache := true
	// if len(values["cache"]) < 0 {
	// 	cache = false
	// }

	doc, err := genenames.GetDoc(geneSymbol)
	if err != nil {
		http.Error(w, "Could not get doc for gene "+geneSymbol, http.StatusInternalServerError)
		return
	}

	summary, err := eutils.GeneSummary(doc.EntrezId)
	if err != nil {
		http.Error(w, "Could not get summary for gene "+geneSymbol, http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(summary)
	if err != nil {
		http.Error(w, "Could not json marshal "+geneSymbol, http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/doc", DocHandler)
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	fmt.Println("Server started on port", port)
	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		fmt.Println(err)
		return
	}
}
