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
	"log"
	"strconv"
)

type User struct {
	Id        string
	FirstName string
	LastName  string
	Age       int
	Weight    float32
}

type TypesService struct {
	RestService `root:"/types-service/" consumes:"application/json" produces:"application/json" realm:"testing"`

	getVarArgs       EndPoint `method:"GET" path:"/var/{...:int}" output:"string" role:"var-user"`
	getVarArgsString EndPoint `method:"GET" path:"/varstring/{...:string}" output:"string"`

	getString            EndPoint `method:"GET" path:"/string/{Bool:bool}/{Int:int}?{flow:int}&{name:string}" output:"string" role:"string-user"`
	getStringSimilarPath EndPoint `method:"GET" path:"/strin?{name:string}" output:"string"`
	getInteger           EndPoint `method:"GET" path:"/int/{Bool:bool}/int/yes/{Int:int}/for" output:"int"`
	getBool              EndPoint `method:"GET" path:"/bool/{Bool:bool}/{Int:int}" output:"bool"`
	getFloat             EndPoint `method:"GET" path:"/float/{Bool:bool}/{Int:int}" output:"float64"`
	getMapInt            EndPoint `method:"GET" path:"/mapint/{Bool:bool}/{Int:int}" output:"map[string]int"`
	getMapStruct         EndPoint `method:"GET" path:"/mapstruct/{Bool:bool}/{Int:int}" output:"map[string]User"`
	getArrayStruct       EndPoint `method:"GET" path:"/arraystruct/{FName:string}/{Age:int}" output:"[]User"`

	postVarArgs EndPoint `method:"POST" path:"/var/{...:int}" postdata:"string"`
	postString  EndPoint `method:"POST" path:"/string/{Bool:bool}/{Int:int}" postdata:"string" role:"post-user"`
	postInteger EndPoint `method:"POST" path:"/int/{Bool:bool}/{Int:int}" postdata:"int" role:"postInt-user"`
	postBool    EndPoint `method:"POST" path:"/bool/{Bool:bool}/{Int:int}" postdata:"bool" `
	postFloat   EndPoint `method:"POST" path:"/float/{Bool:bool}/{Int:int}" postdata:"float64" `

	postMapInt EndPoint `
						method:"POST" 
						path:"/mapint/{Bool:bool}/{Int:int}" 
						postdata:"map[string]int" `

	postMapStruct EndPoint `
						method:"POST" 
						path:"/mapstruct/{Bool:bool}/{Int:int}" 
						postdata:"map[string]User" `

	postArrayStruct EndPoint `
						method:"POST" 
						path:"/arraystruct/{Bool:bool}/{Int:int}" 
						postdata:"[]User"`

	head     EndPoint `method:"HEAD" path:"/bool/{Bool:bool}/{Int:int}"`
	options  EndPoint `method:"OPTIONS" path:"/bool/{Bool:bool}/{Int:int}"`
	doDelete EndPoint `method:"DELETE" path:"/bool/{Bool:bool}/{Int:int}"`
}

type Complex struct {
	Auth       string `Header:""`
	Pathy      int    `Path:"Bool"`
	Query      int    `Query:"flow"`
	CookieUser string `Cookie:"User"`
	CookiePass string `Cookie:"Pass"`
}

var idsInRealm map[string][]string

type TestSessiondata struct {
	id string
}

func (sess *TestSessiondata) SessionId() string {
	return sess.id
}

func (serv TypesService) Head(Bool bool, Int int) {
	rb := serv.ResponseBuilder()
	rb.ETag("12345")
	rb.Age(60 * 30) //30 minutes old

}
func (serv TypesService) DoDelete(Bool bool, Int int) {
	//Will return default response code of 200
}

func (serv TypesService) Options(Bool bool, Int int) {
	rb := serv.ResponseBuilder()
	rb.Allow("GET")
	rb.Allow("HEAD").Allow("POST")
}

func (serv TypesService) GetVarArgs(v ...int) string {
	str := "Start"
	for _, i := range v {
		str += strconv.Itoa(i)
	}
	return str + "End"
}
func (serv TypesService) GetVarArgsString(v ...string) string {
	str := "Start"
	for _, i := range v {
		str += i
	}
	return str + "End"
}
func (serv TypesService) PostVarArgs(name string, varArgs ...int) {
	if name == "hello" && varArgs[0] == 5 && varArgs[1] == 24567 {
		serv.ResponseBuilder().SetResponseCode(200)
	} else {
		serv.ResponseBuilder().SetResponseCode(400)
	}

}
func (serv TypesService) GetStringSimilarPath(name string) string {
	return "Yebo-Yes-" + name
}

func (serv TypesService) GetString(Bool bool, Int int, Flow int, Name string) string {
	return "Hello" + strconv.FormatBool(Bool) + strconv.Itoa(Int) + "/" + Name + strconv.Itoa(Flow)
}
func (serv TypesService) GetInteger(Bool bool, Int int) int {
	return Int - 5
}
func (serv TypesService) GetBool(Bool bool, Int int) bool {
	return Bool
}
func (serv TypesService) GetFloat(Bool bool, Int int) float64 {
	return 111.111 * float64(Int)
}
func (serv TypesService) GetMapInt(Bool bool, Int int) map[string]int {
	mp := make(map[string]int, 0)
	mp["One"] = 1
	mp["Two"] = 2
	mp["Three"] = 3
	return mp
}
func (serv TypesService) GetMapStruct(Bool bool, Int int) map[string]User {
	mp := make(map[string]User, 0)
	mp["One"] = User{"1", "David1", "Gueta1", 35, 123}
	mp["Two"] = User{"2", "David2", "Gueta2", 35, 123}
	mp["Three"] = User{"3", "David3", "Gueta3", 35, 123}
	return mp
}

func (serv TypesService) GetArrayStruct(FName string, Age int) []User {
	users := make([]User, 0)
	users = append(users, User{"user1", FName, "Soap", Age, 89.7})
	users = append(users, User{"user2", FName, "Soap2", Age, 89.7})
	return users
}

func (serv TypesService) PostString(posted string, Bool bool, Int int) {
	if posted == "Hello" && Bool && Int == 5 {
		serv.ResponseBuilder().SetResponseCode(200)
	} else {
		serv.ResponseBuilder().SetResponseCode(400)
	}
	log.Println("posted:", posted)
}
func (serv TypesService) PostInteger(posted int, Bool bool, Int int) {
	if posted == 6 && Bool && Int == 5 {
		serv.ResponseBuilder().SetResponseCode(200)
	} else {
		serv.ResponseBuilder().SetResponseCode(400)
	}
	log.Println("posted:", posted)
}
func (serv TypesService) PostBool(posted bool, Bool bool, Int int) {
	if !posted && Bool && Int == 5 {
		serv.ResponseBuilder().SetResponseCode(200)
	} else {
		serv.ResponseBuilder().SetResponseCode(400)
	}
	log.Println("posted:", posted)
}
func (serv TypesService) PostFloat(posted float64, Bool bool, Int int) {
	if posted == 34.56788 && Bool && Int == 5 {
		serv.ResponseBuilder().SetResponseCode(200)
	} else {
		serv.ResponseBuilder().SetResponseCode(400)
	}
	log.Println("posted:", posted)
}
func (serv TypesService) PostMapInt(posted map[string]int, Bool bool, Int int) {

	if posted["One"] == 111 && posted["Two"] == 222 && Bool && Int == 5 {
		serv.ResponseBuilder().SetResponseCode(200)
	} else {
		serv.ResponseBuilder().SetResponseCode(400)
	}
	log.Println("posted map One:", posted["One"])
	log.Println("posted map Two:", posted["Two"])
}
func (serv TypesService) PostMapStruct(posted map[string]User, Bool bool, Int int) {
	rb := serv.ResponseBuilder()
	if posted["One"].FirstName == "David1" && posted["Two"].LastName == "Gueta2" && Bool && Int == 5 {
		rb.SetResponseCode(200)
	} else {
		rb.SetResponseCode(400)
	}

	rb.Write([]byte(posted["One"].FirstName + posted["One"].LastName + posted["One"].Id))
	rb.Write([]byte(posted["Two"].FirstName + posted["Two"].LastName + posted["Two"].Id))

	log.Println("posted map One:", posted["One"].FirstName, posted["One"].LastName, posted["One"].Id)
	log.Println("posted map Two:", posted["Two"].FirstName, posted["Two"].LastName, posted["Two"].Id)
}
func (serv TypesService) PostArrayStruct(posted []User, Bool bool, Int int) {
	if posted[0].FirstName == "Joe" && posted[1].LastName == "Soap2" && Bool && Int == 5 {
		serv.ResponseBuilder().SetResponseCode(200)
	} else {
		serv.ResponseBuilder().SetResponseCode(400)
	}
	log.Println("posted Array One:", posted[0].FirstName, posted[0].LastName, posted[0].Id)
	log.Println("posted Array Two:", posted[1].FirstName, posted[1].LastName, posted[1].Id)

}
