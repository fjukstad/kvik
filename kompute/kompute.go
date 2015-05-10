package kompute

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Kompute struct {
	Addr     string
	Username string
	Password string
}

func NewKompute(addr, username, password string) *Kompute {
	komp := new(Kompute)
	komp.Addr = "http://" + addr
	komp.Username = username
	komp.Password = password
	return komp
}

// Plots and stores to file
func (k *Kompute) Plot(fun, args, filetype, filename string) error {
	s, err := k.Call(fun, args)
	if err != nil {
		return err
	}

	url := s.graphics + "/" + filetype
	resp, err := http.Get(url)

	if err != nil {
		return nil
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	plot, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, plot, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Execcutes the command and returns the results formatted in the specified
// format e.g. json or csv
func (k *Kompute) Rpc(fun, args, format string) (string, error) {

	fmt.Println("-H \"Content-Type: application/json\" -d '", args, "'")

	s, err := k.Call(fun, args)

	if err != nil {
		fmt.Println("Call error != nil")
		return "", err
	}

	s.GetResult(k, format)

	return s.Result, err
}

func (k *Kompute) Call(fun, args string) (s *Session, err error) {

	url := k.getUrl(fun)
	postArgs := strings.NewReader(args)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, postArgs)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	header := map[string][]string{
		"Content-Type": {"application/json"},
	}

	req.Header = header
	req.SetBasicAuth(k.Username, k.Password)

	resp, err := client.Do(req)

	//resp, err := http.Post(url, "application/json", postArgs)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("empty body")
		return nil, err
	}

	if resp.StatusCode != 201 {
		fmt.Println("Statuscode != 201")
		return nil, errors.New(string(body))
	}

	s = new(Session)

	s.Key = resp.Header.Get("X-Ocpu-Session")
	s.Url = resp.Header.Get("Location")

	output := strings.Split(string(body), "\n")
	for _, line := range output {
		switch {
		case strings.Contains(line, ".val"):
			s.val = k.Addr + line
		case strings.Contains(line, "stdout"):
			s.stdout = k.Addr + line
		case strings.Contains(line, "source"):
			s.source = k.Addr + line
		case strings.Contains(line, "info"):
			s.info = k.Addr + line
		case strings.Contains(line, "DESCRIPTION"):
			s.description = k.Addr + line
		case strings.Contains(line, "console"):
			s.console = k.Addr + line
		case strings.Contains(line, "graphics"):
			s.graphics = k.Addr + line
		}
	}

	return s, nil

}

func (k *Kompute) getUrl(fun string) string {
	return k.Addr + "/ocpu/library/" + fun
}

type Session struct {
	Key         string
	Url         string
	Result      string
	console     string
	stdout      string
	val         string
	info        string
	description string
	graphics    string
	source      string
}

func (s *Session) GetResult(k *Kompute, format string) (string, error) {
	url := s.val + "/" + format

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.SetBasicAuth(k.Username, k.Password)

	resp, err := client.Do(req)

	//resp, err := http.Get(url)

	if err != nil {
		return "", nil
	}

	if resp.StatusCode != 200 {
		fmt.Println(url)
		return "", errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("bodu thing wrong", err)
		return "", nil
	}

	s.Result = string(body)

	return s.Result, nil
}
