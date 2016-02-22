package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/fjukstad/r"
)

func main() {
	t0 := time.Now()
	url := "http://localhost:8181/call"

	c := R.RCall{"stats", "rnorm", "n=10000000"}

	b, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	postArgs := strings.NewReader(string(b))

	resp, err := http.Post(url, "application/json", postArgs)
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err = http.Get("http://localhost:8181/get/" + string(body) + "/json")
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))

}
