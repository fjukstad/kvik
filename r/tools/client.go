package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fjukstad/r"
)

func UploadPackage(url, src string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	file.Close()
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return err
	}
	part.Write(fileContents)
	err = writer.Close()
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Response:", string(responseBody))
	return nil
}

func TestFunctionality(server, pkg, fun, args string) (string, error) {

	url := "http://" + server + "/call"

	c := R.RCall{pkg, fun, args}

	b, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	postArgs := strings.NewReader(string(b))

	resp, err := http.Post(url, "application/json", postArgs)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	resp, err = http.Get("http://" + server + "/get/" + string(body) + "/json")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(body), nil

}

func main() {
	t0 := time.Now()
	var server = flag.String("server", "localhost:8181", "ip:port of server")
	flag.Parse()

	out, err := TestFunctionality(*server, "stats", "rnorm", "n=10")
	fmt.Println(out, err)

	url := "http://" + *server + "/install"
	err = UploadPackage(url, "packages/addman_0.1.tgz")
	if err != nil {
		fmt.Println(err)
		return
	}

	out, err = TestFunctionality(*server, "addman", "hello", "")
	fmt.Println(out, err)
	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))

}
