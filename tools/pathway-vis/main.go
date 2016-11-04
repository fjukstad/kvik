package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fjukstad/kvik/kegg"
	"github.com/gorilla/mux"
)

func PathwayHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	graph := kegg.PathwayGraph(id)
	b, err := json.Marshal(graph)

	if err != nil {
		fmt.Println("Could not marshal response ", err)
		return
	}

	w.Write(b)

}

func FileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("FILEMAN")
	vars := mux.Vars(r)
	filename := vars["filename"]

	if filename == "" {
		filename = "index.html"
	}

	fmt.Println(filename)
	http.ServeFile(w, r, filename)
	return
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/pathway/{id}", PathwayHandler)
	http.Handle("/", http.FileServer(http.Dir("./")))

	http.Handle("/pathway/", r)

	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)

}
