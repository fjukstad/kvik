package r

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	Addr     string
	Username string
	Password string
}

func (c *Client) Call(pkg, fun, args string) (string, error) {

	url := "http://" + c.Addr + "/call"
	b, _ := json.Marshal(Call{pkg, fun, args})
	postArgs := strings.NewReader(string(b))

	resp, err := http.Post(url, "application/json", postArgs)
	if err != nil {
		fmt.Println("Could not post to R server", err)
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Could not read body", body, err)
		return "", err
	}

	return string(body), nil
}

func (c *Client) Get(key, format string) ([]byte, error) {

	url := "http://" + c.Addr + "/get/" + key + "/" + format
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
