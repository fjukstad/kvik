//Copyright 2011 Siyabonga Dlamini (siyabonga.dlamini@gmail.com). All rights reserved.
//
//Redistribution and use in source and binary forms, with or without
//modification, are permitted provided that the following conditions
//are met:
//
//  1. Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer.
//
//  2. Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer
//     in the documentation and/or other materials provided with the
//     distribution.
//
//THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR
//IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
//OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
//IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
//PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS;
//OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
//WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
//OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
//ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package gorest

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var sharedClient *http.Client

//Use this if you have a *http.Client instance that you specifically want to use. 
//Otherwise just use NewRequestBuilder(), which uses the http.Client maintained by GoRest.
func NewRequestBuilderFromClient(client *http.Client, url string) (*RequestBuilder, error) {
	req, err := http.NewRequest(GET, url, nil)
	if err != nil {
		return nil, err
	}
	rb := RequestBuilder{client, Application_Json, req}
	return &rb, nil
}

//This creates a new RequestBuilder, backed by GoRest's internally managed http.Client.
//Although http.Client is useable concurrently, an instance of RequestBuilder is not safe for this. 
//Because of http.Client's persistent(cached TCP connections) and concurrent nature, 
//this can be used safely multiple times from different go routines. 
func NewRequestBuilder(url string) (*RequestBuilder, error) {
	if sharedClient == nil {
		sharedClient = new(http.Client) //DefaultClient
	}
	req, err := http.NewRequest(GET, url, nil)
	if err != nil {
		return nil, err
	}
	rb := RequestBuilder{sharedClient, Application_Json, req}
	return &rb, nil
}

type RequestBuilder struct {
	client             *http.Client
	defaultContentType string
	_req               *http.Request
}

func (this *RequestBuilder) Request() *http.Request {
	return this._req
}
func (this *RequestBuilder) UseContentType(mime string) *RequestBuilder {
	this.defaultContentType = mime
	return this
}

func (this *RequestBuilder) CacheNoCache() *RequestBuilder {
	this.setCache("no-cache")
	return this
}
func (this *RequestBuilder) CacheNoStore() *RequestBuilder {
	this.setCache("no-store")
	return this
}
func (this *RequestBuilder) CacheMaxAge(seconds int) *RequestBuilder {
	this.setCache("max-age = " + strconv.Itoa(seconds))
	return this
}
func (this *RequestBuilder) CacheStale(seconds int) *RequestBuilder {
	this.setCache("max-stale = " + strconv.Itoa(seconds))
	return this
}
func (this *RequestBuilder) CacheMinFresh(seconds int) *RequestBuilder {
	this.setCache("min-fresh = " + strconv.Itoa(seconds))
	return this
}
func (this *RequestBuilder) CacheOnlyIfCached() *RequestBuilder {
	this.setCache("only-if-cached")
	return this
}
func (this *RequestBuilder) CacheClearAllOptions() *RequestBuilder {
	this._req.Header.Del("Cache-control")
	return this
}
func (this *RequestBuilder) setCache(option string) {
	this._req.Header.Add("Cache-control", option)
}

func (this *RequestBuilder) Accept(mime string) *RequestBuilder {
	this._req.Header.Add("Accept", mime)
	return this
}
func (this *RequestBuilder) AcceptClear() *RequestBuilder {
	this._req.Header.Del("Accept")
	return this
}
func (this *RequestBuilder) AcceptCharSet(set string) *RequestBuilder {
	this._req.Header.Add("Accept-Charset", set)
	return this
}
func (this *RequestBuilder) AcceptCharSetClear() *RequestBuilder {
	this._req.Header.Del("Accept-Charset")
	return this
}
func (this *RequestBuilder) AcceptEncoding(option string) *RequestBuilder {
	this._req.Header.Add("Accept-Encoding", option)
	return this
}
func (this *RequestBuilder) AcceptEncodingClear() *RequestBuilder {
	this._req.Header.Del("Accept-Encoding")
	return this
}
func (this *RequestBuilder) AcceptLanguage(lang string) *RequestBuilder {
	this._req.Header.Add("Accept-Language", lang)
	return this
}
func (this *RequestBuilder) AcceptLanguageClear() *RequestBuilder {
	this._req.Header.Del("Accept-Language")
	return this
}

func (this *RequestBuilder) ConnectionKeepAlive() *RequestBuilder {
	this._req.Header.Set("Connection", "keep-alive")
	return this
}
func (this *RequestBuilder) ConnectionClose() *RequestBuilder {
	this._req.Header.Set("Connection", "close")
	return this
}
func (this *RequestBuilder) AddCookie(cookie *http.Cookie) *RequestBuilder {
	this._req.AddCookie(cookie)
	return this
}

func (this *RequestBuilder) Delete() (*http.Response, error) {
	//	u, err := url.Parse(url_)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	//this._req.URL = u
	this._req.Method = DELETE

	return this.client.Do(this._req)
}

func (this *RequestBuilder) Head() (*http.Response, error) {
	this._req.Method = HEAD
	return this.client.Do(this._req)
}

func (this *RequestBuilder) Options(opts *[]string) (*http.Response, error) {
	this._req.Method = OPTIONS

	res, err := this.client.Do(this._req)
	if err != nil {
		return res, err
	}

	for _, str := range res.Header["Allow"] {
		*opts = append(*opts, strings.Trim(str, " "))
	}
	return res, err
}

func (this *RequestBuilder) Get(i interface{}, expecting int) (*http.Response, error) {
	//this._req.URL = u
	this._req.Method = GET

	res, err := this.client.Do(this._req)
	if err != nil {
		return res, err
	}

	if res.StatusCode == expecting {
		buf := new(bytes.Buffer)
		io.Copy(buf, res.Body)
		res.Body.Close()
		err = BytesToInterface(buf, i, this.defaultContentType)
		return res, nil
	}

	return res, errors.New(res.Status)
}
func (this *RequestBuilder) Post(i interface{}) (*http.Response, error) {
	this._req.Method = POST
	bb, err := InterfaceToBytes(i, this.defaultContentType)
	if err != nil {
		return nil, err
	}
	this._req.Body = ioutil.NopCloser(bytes.NewBuffer(bb))

	return this.client.Do(this._req)

}
