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
	"testing"
)

func testPaths(t *testing.T) {
	deleteMixed1(t)
	deleteMixed2(t)
	optionsMixed(t)
	getMixed(t)
}

func deleteMixed1(t *testing.T) {
	//deleteMixed1 EndPoint `method:"DELETE" 	path:"/bool/{Bool:bool}/mix1/{Int:int}"`
	//*******************************

	rb, _ := NewRequestBuilder(RootPath + "paths-service/bool/true/mix1/5" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Delete()
	AssertEqual(res.StatusCode, 200, "Delete ResponseCode", t)
}
func deleteMixed2(t *testing.T) {
	//deleteMixed2 EndPoint `method:"DELETE" 	path:"/bool/{Bool:bool}/mix2/{Int:int}"`
	//*******************************

	rb, _ := NewRequestBuilder(RootPath + "paths-service/bool/true/mix2/5" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Delete()
	AssertEqual(res.StatusCode, 200, "Delete ResponseCode", t)
}

func optionsMixed(t *testing.T) {
	//optionsMixed EndPoint `method:"OPTIONS" path:"/bool/{Bool:bool}/mix1/{Int:int}"`
	//*******************************

	strArr := make([]string, 0)
	rb, _ := NewRequestBuilder(RootPath + "paths-service/bool/true/mix1/5" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Options(&strArr)
	AssertEqual(res.StatusCode, 200, "Options ResponseCode", t)
	AssertEqual(len(strArr), 3, "Options - slice length", t)
	if len(strArr) == 3 {
		AssertEqual(strArr[0], GET, "Options", t)
		AssertEqual(strArr[1], HEAD, "Options", t)
		AssertEqual(strArr[2], POST, "Options", t)
	}
}

func getMixed(t *testing.T) {
	//getMixed     EndPoint `method:"GET" 	path:"/bool/{Bool:bool}/mix1/{Int:int}"`
	//*******************************	

	rb, _ := NewRequestBuilder(RootPath + "paths-service/bool/true/mix1/5" + xrefStr)
	rb.AddCookie(cook)
	//GET string

	res, _ := rb.Get(&str, 200)
	AssertEqual(res.StatusCode, 200, "Get string ResponseCode", t)
	AssertEqual(str, "Hello", "Get getMixed", t)

}
