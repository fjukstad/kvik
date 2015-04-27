package kompute

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type Kompute struct {
	Addr string
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

	s, err := k.Call(fun, args)

	if err != nil {
		return "", err
	}

	s.GetResult(format)

	return s.Result, err
}

func (k *Kompute) Call(fun, args string) (s *Session, err error) {

	url := getUrl(k.Addr, fun)
	postArgs := strings.NewReader(args)

	resp, err := http.Post(url, "application/json", postArgs)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 201 {
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

func getUrl(hostname, fun string) string {
	return hostname + "/ocpu/library/" + fun
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

func (s *Session) GetResult(format string) (string, error) {
	url := s.val + "/" + format
	resp, err := http.Get(url)

	if err != nil {
		return "", nil
	}

	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	s.Result = string(body)

	return s.Result, nil
}
