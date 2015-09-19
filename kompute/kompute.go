package kompute

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	//"time"
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

	url := s.Graphics + "/" + filetype
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

	//	fmt.Println("-H \"Content-Type: application/json\" -d '", args, "'")

	s, err := k.Call(fun, args)

	if err != nil {

		fmt.Println("Call error != nil", err)
		s, err = k.Call(fun, args)
		if err != nil {
			fmt.Println("failed a second time...", err)
			return "", err
		}

	}

	//	time.Sleep(100 * time.Millisecond)

	res, err := s.GetResult(k, format)
	if err != nil {
		res, err = s.GetResult(k, format)
		fmt.Println("Get result error second time ")
	}
	return res, err
}

func (k *Kompute) Call(fun, args string) (s *Session, err error) {

	url := k.getUrl(fun)
	postArgs := strings.NewReader(args)

	fmt.Println(url, postArgs)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, postArgs)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var contentType string
	if strings.Contains(args, "{") {
		contentType = "application/json"
	} else {
		contentType = "application/x-www-form-urlencoded"
	}

	header := map[string][]string{
		"Content-Type": {contentType},
	}

	req.Header = header
	req.SetBasicAuth(k.Username, k.Password)

	defer req.Body.Close()

	var resp *http.Response
	resp, err = client.Do(req)

	maxretry := 100
	//resp, err := http.Post(url, "application/json", postArgs)
	if err != nil {
		r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
		for i := 0; i < maxretry; i++ {
			randTime := r.Intn(200)
			time.Sleep(time.Duration(randTime) * time.Millisecond)
			resp, err = client.Do(req)
			if err == nil {
				break
			}
		}
		if err != nil {
			return nil, err
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("empty body")
		return nil, err
	}

	defer resp.Body.Close()

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
			s.Graphics = k.Addr + line
		}
	}

	return s, nil

}

func (k *Kompute) getUrl(fun string) string {
	return k.Addr + "/ocpu/library/" + fun
}

func (k *Kompute) Get(key, filetype string) ([]byte, error) {
	var url string
	if strings.Contains(filetype, "png") || strings.Contains(filetype, "pdf") {
		url = k.Addr + "/ocpu/tmp/" + key + "/graphics/last/" + filetype
	} else {
		url = k.Addr + "/ocpu/tmp/" + key + "/R/.val/" + filetype
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Culd not create request", err)
		return nil, err
	}

	req.SetBasicAuth(k.Username, k.Password)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Could not get", err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
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
	Graphics    string
	source      string
}

func (s *Session) GetUrl(format string) (url string) {
	return s.val + "/" + format
}

func (s *Session) GetResult(k *Kompute, format string) (string, error) {
	url := s.GetUrl(format)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("failed second time")
			return "", err
		}
	}
	header := map[string][]string{
		"Content-Type": {"application/json"},
	}

	req.Header = header
	req.SetBasicAuth(k.Username, k.Password)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Client did not")
		return "", err
	}

	if resp.StatusCode != 200 {
		fmt.Println(url)
		fmt.Println(req)
		fmt.Println(resp)
		error, _ := ioutil.ReadAll(resp.Body)
		errorText := string(error)
		fmt.Println("Status code not 200 in GetResult", string(error))
		resp, err = client.Do(req)
		if err != nil && resp.StatusCode != 200 {
			fmt.Println("status code second time jesus christ")
			return "", errors.New(errorText)
		}

	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("bodu thing wrong", err)
		return "", nil
	}

	s.Result = string(body)

	return s.Result, nil
}
