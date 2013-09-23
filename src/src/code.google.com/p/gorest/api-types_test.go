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
	"net/http"
	"strings"
	"testing"
)

var cook = new(http.Cookie)
var xrefStr = "?xsrft=12345"
var str = "Hell"

func testTypes(t *testing.T) {

	cook.Name = "X-Xsrf-Cookie"
	cook.Value = "12345"

	getString(t)
	getStringSimilarPath(t)
	getVarArgsString(t)
	getInteger(t)
	getBool(t)
	getFloat(t)
	getMapInt(t)
	getMapStruct(t)
	getArrayStruct(t)

	postVarArgs(t)
	postVarArgs(t)
	postString(t)
	postInteger(t)
	postBool(t)
	postFloat(t)
	postMapInt(t)
	postMapStruct(t)
	postArrayStruct(t)

	head(t)
	options(t)
	doDelete(t)
	doDeleteMixed1(t)
	doDeleteMixed2(t)
}

func getVarArgs(t *testing.T) {
	//getVarArgs       EndPoint `method:"GET" path:"/var/{...:int}" output:"string" role:"var-user"`
	//*******************************	

	rb, _ := NewRequestBuilder(RootPath + "types-service/var/1/2/3/4/5/6/7/8" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Get(&str, 200)
	AssertEqual(res.StatusCode, 200, "Get var-args Int ResponseCode", t)
	AssertEqual(str, "Start12345678End", "Get var-args Int", t)
}

func postVarArgs(t *testing.T) {
	//postVarArgs      EndPoint `method:"POST" path:"/var/{...:int}" postdata:"string"`
	//*******************************	

	rb, _ := NewRequestBuilder(RootPath + "types-service/var/5/24567" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Post("hello")
	AssertEqual(res.StatusCode, 200, "Post Var args", t)
}

func getVarArgsString(t *testing.T) {
	//getVarArgsString EndPoint `method:"GET" path:"/varstring/{...:string}" output:"string"`
	//*******************************	

	rb, _ := NewRequestBuilder(RootPath + "types-service/varstring/One/Two/Three" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Get(&str, 200)
	AssertEqual(res.StatusCode, 200, "Get var-args string ResponseCode", t)
	AssertEqual(str, "StartOneTwoThreeEnd", "Get var-args string", t)
}

func getString(t *testing.T) {
	//getString            EndPoint `method:"GET" path:"/string/{Bool:bool}/{Int:int}?{flow:int}&{name:string}" output:"string" role:"string-user"`
	//*******************************	

	rb, _ := NewRequestBuilder(RootPath + "types-service/string/true/5" + xrefStr + "&name=Nameed&flow=6")
	rb.AddCookie(cook)
	//GET string

	res, _ := rb.Get(&str, 200)
	AssertEqual(res.StatusCode, 200, "Get string ResponseCode", t)
	AssertEqual(str, "Hellotrue5/Nameed6", "Get string", t)

}

func getStringSimilarPath(t *testing.T) {
	//getStringSimilarPath EndPoint `method:"GET" path:"/strin?{name:string}" output:"string"`
	//*******************************	

	rb, _ := NewRequestBuilder(RootPath + "types-service/strin" + xrefStr + "&name=Nameed")
	rb.AddCookie(cook)
	res, _ := rb.Get(&str, 200)
	AssertEqual(res.StatusCode, 200, "Get string ResponseCode", t)
	AssertEqual(str, "Yebo-Yes-Nameed", "Get string similar path", t)

	rb, _ = NewRequestBuilder(RootPath + "types-service/string/true/5" + xrefStr + "&name=Nameed")
	rb.AddCookie(cook)
	res, _ = rb.Get(&str, 200)
	AssertEqual(res.StatusCode, 200, "Get string ResponseCode", t)
	AssertEqual(str, "Hellotrue5/Nameed0", "Get string", t)

	rb, _ = NewRequestBuilder(RootPath + "types-service/string/true/5" + xrefStr + "&flow=6")
	rb.AddCookie(cook)
	res, _ = rb.Get(&str, 200)
	AssertEqual(res.StatusCode, 200, "Get string ResponseCode", t)
	AssertEqual(str, "Hellotrue5/6", "Get string", t)

	rb, _ = NewRequestBuilder(RootPath + "types-service/string/true/5" + xrefStr + "&flow=")
	rb.AddCookie(cook)
	res, _ = rb.Get(&str, 200)
	AssertEqual(res.StatusCode, 200, "Get string ResponseCode", t)
	AssertEqual(str, "Hellotrue5/0", "Get string", t)

	rb, _ = NewRequestBuilder(RootPath + "types-service/string/true/5" + xrefStr + "&flow")
	rb.AddCookie(cook)
	res, _ = rb.Get(&str, 200)
	AssertEqual(res.StatusCode, 200, "Get string ResponseCode", t)
	AssertEqual(str, "Hellotrue5/0", "Get string", t)
}

func getInteger(t *testing.T) {
	//getInteger           EndPoint `method:"GET" path:"/int/{Bool:bool}/int/yes/{Int:int}/for" output:"int"`
	//*******************************	

	inter := -2
	rb, _ := NewRequestBuilder(RootPath + "types-service/int/true/int/yes/2/for" + xrefStr + "&name=Nameed&flow=6")
	rb.AddCookie(cook)
	res, _ := rb.Get(&inter, 200) //The query aurguments here just to be ignored
	AssertEqual(res.StatusCode, 200, "Get int ResponseCode", t)
	AssertEqual(inter, -3, "Get int", t)
}

func getBool(t *testing.T) {
	//getBool              EndPoint `method:"GET" path:"/bool/{Bool:bool}/{Int:int}" output:"bool"`
	//*******************************	
	bl := true
	rb, _ := NewRequestBuilder(RootPath + "types-service/bool/false/2" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Get(&bl, 200)
	AssertEqual(res.StatusCode, 200, "Get int ResponseCode", t)
	AssertEqual(bl, false, "Get Bool", t)
}

func getFloat(t *testing.T) {
	//getFloat             EndPoint `method:"GET" path:"/float/{Bool:bool}/{Int:int}" output:"float64"`
	//*******************************

	fl := 2.4
	rb, _ := NewRequestBuilder(RootPath + "types-service/float/false/2" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Get(&fl, 200)
	AssertEqual(res.StatusCode, 200, "Get Float ResponseCode", t)
	AssertEqual(fl, 222.222, "Get Float", t)
}

func getMapInt(t *testing.T) {
	//getMapInt            EndPoint `method:"GET" path:"/mapint/{Bool:bool}/{Int:int}" output:"map[string]int"`
	//*******************************

	mp := make(map[string]int)
	rb, _ := NewRequestBuilder(RootPath + "types-service/mapint/false/2" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Get(&mp, 200)
	AssertEqual(res.StatusCode, 200, "Get Float ResponseCode", t)
	AssertEqual(mp["One"], 1, "Get Map Int", t)
	AssertEqual(mp["Two"], 2, "Get Map Int", t)
}

func getMapStruct(t *testing.T) {
	//getMapStruct         EndPoint `method:"GET" path:"/mapstruct/{Bool:bool}/{Int:int}" output:"map[string]User"`
	//*******************************

	mpu := make(map[string]User)
	rb, _ := NewRequestBuilder(RootPath + "types-service/mapstruct/false/2" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Get(&mpu, 200)
	AssertEqual(res.StatusCode, 200, "Get Map struct ResponseCode", t)
	AssertEqual(mpu["One"].Id, "1", "Get Map struct", t)
	AssertEqual(mpu["Two"].Id, "2", "Get Map struct", t)
	AssertEqual(mpu["Two"].FirstName, "David2", "Get Map struct", t)
	AssertEqual(mpu["Two"].LastName, "Gueta2", "Get Map struct", t)
}

func getArrayStruct(t *testing.T) {
	//getArrayStruct       EndPoint `method:"GET" path:"/arraystruct/{FName:string}/{Age:int}" output:"[]User"`
	//*******************************
	au := make([]User, 0)
	rb, _ := NewRequestBuilder(RootPath + "types-service/arraystruct/Sandy/2" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Get(&au, 200)
	AssertEqual(res.StatusCode, 200, "Get Array struct ResponseCode", t)
	if res.StatusCode == 200 {
		AssertEqual(au[0].Id, "user1", "Get Array Struct", t)
		AssertEqual(au[0].FirstName, "Sandy", "Get Array Struct", t)
	}
}

func postString(t *testing.T) {
	//postString      EndPoint `method:"POST" path:"/string/{Bool:bool}/{Int:int}" postdata:"string" role:"post-user"`
	//*******************************
	rb, _ := NewRequestBuilder(RootPath + "types-service/string/true/5" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Post("Hello")
	AssertEqual(res.StatusCode, 200, "Post String", t)
}

func postInteger(t *testing.T) {
	//postInteger     EndPoint `method:"POST" path:"/int/{Bool:bool}/{Int:int}" postdata:"int" role:"postInt-user"`
	//*******************************

	//POST Int requires the postInt-user role, which only user fox has
	rb, _ := NewRequestBuilder(RootPath + "types-service/int/true/5" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Post(6)
	AssertEqual(res.StatusCode, 403, "Post Integer wrong user", t)

	cook2 := new(http.Cookie)
	cook2.Name = "X-Xsrf-Cookie"
	cook2.Value = "fox"

	xrefStr2 := "?xsrft=fox"

	rb2, _ := NewRequestBuilder(RootPath + "types-service/int/true/5" + xrefStr2)
	rb2.AddCookie(cook2)
	res, _ = rb2.Post(6)
	AssertEqual(res.StatusCode, 200, "Post Integer correct user", t)
}

func postBool(t *testing.T) {
	//postBool        EndPoint `method:"POST" path:"/bool/{Bool:bool}/{Int:int}" postdata:"bool" `
	//*******************************

	rb, _ := NewRequestBuilder(RootPath + "types-service/bool/true/5" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Post(false)
	AssertEqual(res.StatusCode, 200, "Post Boolean", t)
}

func postFloat(t *testing.T) {
	//postFloat       EndPoint `method:"POST" path:"/float/{Bool:bool}/{Int:int}" postdata:"float64" `
	//*******************************

	rb, _ := NewRequestBuilder(RootPath + "types-service/float/true/5" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Post(34.56788)
	AssertEqual(res.StatusCode, 200, "Post Float", t)
}

func postMapInt(t *testing.T) {
	//postMapInt      EndPoint `method:"POST" path:"/mapint/{Bool:bool}/{Int:int}" postdata:"map[string]int" `
	//*******************************

	mi := make(map[string]int, 0)
	mi["One"] = 111
	mi["Two"] = 222
	rb, _ := NewRequestBuilder(RootPath + "types-service/mapint/true/5" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Post(mi)
	AssertEqual(res.StatusCode, 200, "Post Integer Map", t)
}

func postMapStruct(t *testing.T) {
	//postMapStruct   EndPoint `method:"POST" path:"/mapstruct/{Bool:bool}/{Int:int}" postdata:"map[string]User" `
	//*******************************

	mu := make(map[string]User, 0)
	mu["One"] = User{"111", "David1", "Gueta1", 35, 123}
	mu["Two"] = User{"222", "David2", "Gueta2", 35, 123}
	rb, _ := NewRequestBuilder(RootPath + "types-service/mapstruct/true/5" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Post(mu)
	AssertEqual(res.StatusCode, 200, "Post Struct Map", t)
	//AssertEqual(res., 200, "Post Struct Map", t)
}

func postArrayStruct(t *testing.T) {
	//postArrayStruct EndPoint `method:"POST" path:"/arraystruct/{Bool:bool}/{Int:int}" postdata:"[]User"`
	//*******************************

	users := make([]User, 0)
	users = append(users, User{"user1", "Joe", "Soap", 19, 89.7})
	users = append(users, User{"user2", "Jose", "Soap2", 15, 89.7})

	rb, _ := NewRequestBuilder(RootPath + "types-service/arraystruct/true/5" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Post(users)
	AssertEqual(res.StatusCode, 200, "Post Struct Array", t)
}

func head(t *testing.T) {
	//head     EndPoint `method:"HEAD" path:"/bool/{Bool:bool}/{Int:int}"`
	//*******************************

	rb, _ := NewRequestBuilder(RootPath + "types-service/bool/false/2" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Head()
	AssertEqual(res.StatusCode, 200, "Head ResponseCode", t)
	AssertEqual(res.Header.Get("ETag"), "12345", "Head Header ETag", t)
	AssertEqual(len(res.Header["Age"]), 1, "Head Header Age - slice length", t)
	if len(res.Header["Age"]) == 1 {
		AssertEqual(strings.Trim(res.Header["Age"][0], " "), "1800", "Head Header Age", t)
	}
}

func options(t *testing.T) {
	//options  EndPoint `method:"OPTIONS" path:"/bool/{Bool:bool}/{Int:int}"`
	//*******************************

	strArr := make([]string, 0)
	rb, _ := NewRequestBuilder(RootPath + "types-service/bool/false/2" + xrefStr)
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

func doDelete(t *testing.T) {
	//doDelete EndPoint `method:"DELETE" path:"/bool/{Bool:bool}/{Int:int}"`
	//*******************************

	rb, _ := NewRequestBuilder(RootPath + "types-service/bool/false/2" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Delete()
	AssertEqual(res.StatusCode, 200, "Delete ResponseCode", t)
}

func doDeleteMixed1(t *testing.T) {
	//doDeleteMixed1 EndPoint `method:"DELETE" path:"/bool/{Bool:bool}/mix1/{Int:int}"`
	//*******************************

	rb, _ := NewRequestBuilder(RootPath + "types-service/bool/false/2" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Delete()
	AssertEqual(res.StatusCode, 200, "Delete ResponseCode", t)
}
func doDeleteMixed2(t *testing.T) {
	//doDeleteMixed2 EndPoint `method:"DELETE" path:"/bool/{Bool:bool}/mix2/{Int:int}"`
	//*******************************

	rb, _ := NewRequestBuilder(RootPath + "types-service/bool/false/2" + xrefStr)
	rb.AddCookie(cook)
	res, _ := rb.Delete()
	AssertEqual(res.StatusCode, 200, "Delete ResponseCode", t)
}
