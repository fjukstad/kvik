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
	"encoding/json"
	"strconv"
	"testing"
)

func TestInterfaceToBytes(t *testing.T) {

	bytes, _ := InterfaceToBytes(12345, "application/json")
	AssertEqual(string(bytes), "12345", "Integer marshall", t)

	bytes, _ = InterfaceToBytes("Hello", "application/json")
	AssertEqual(string(bytes), "Hello", "String marshall", t)

	bytes, _ = InterfaceToBytes(true, "application/json")
	AssertEqual(string(bytes), "true", "Bool marshall", t)

	bytes, _ = InterfaceToBytes(36.6, "application/json")
	AssertEqual(string(bytes), "36.6", "Float marshall", t)

	bytes, _ = InterfaceToBytes(-37, "application/json")
	AssertEqual(string(bytes), "-37", "Uint marshall", t)

	u := new(User)
	u.FirstName = "David"
	u.LastName = "Coperfield"
	u.Age = 20

	bytes, _ = InterfaceToBytes(u, "application/json")
	AssertEqual(string(bytes), `{"Id":"","FirstName":"David","LastName":"Coperfield","Age":20,"Weight":0}`, "Struct marshall", t)

	userArr := make([]User, 0)
	u2 := *u
	u2.Age = 30
	userArr = append(userArr, *u)
	userArr = append(userArr, u2)
	bytes, _ = InterfaceToBytes(userArr, "application/json")
	AssertEqual(string(bytes), `[{"Id":"","FirstName":"David","LastName":"Coperfield","Age":20,"Weight":0},{"Id":"","FirstName":"David","LastName":"Coperfield","Age":30,"Weight":0}]`, "Array marshall", t)

	userMap := make(map[string]User, 0)
	userMap["One"] = *u
	userMap["Two"] = u2

	bytes, _ = InterfaceToBytes(userMap, "application/json")
	AssertEqual(string(bytes), `{"One":{"Id":"","FirstName":"David","LastName":"Coperfield","Age":20,"Weight":0},"Two":{"Id":"","FirstName":"David","LastName":"Coperfield","Age":30,"Weight":0}}`, "Map marshall", t)

}

func TestBytesToI(t *testing.T) {
	bully := true
	i := 0
	ui := -1
	str := ""
	fl := 34.5

	byt := bytes.NewBufferString("36")
	BytesToInterface(byt, &i, "")
	AssertEqual(i, 36, "Integer unmarshall", t)

	byt = bytes.NewBufferString("false")
	BytesToInterface(byt, &bully, "")
	AssertEqual(bully, false, "Bool unmarshall", t)

	byt = bytes.NewBufferString("-12")
	BytesToInterface(byt, &ui, "")
	AssertEqual(ui, -12, "UInt unmarshall", t)

	byt = bytes.NewBufferString("36.7787")
	BytesToInterface(byt, &fl, "")
	AssertEqual(fl, 36.7787, "Float unmarshall", t)

	byt = bytes.NewBufferString("Hello")
	BytesToInterface(byt, &str, "")
	AssertEqual(str, "Hello", "String unmarshall", t)

	u := new(User)
	u.FirstName = "David"
	u.LastName = "Coperfield"
	u.Age = 20

	by, _ := json.Marshal(u)

	//Try single user
	byt = bytes.NewBuffer(by)
	//println("User",string(byt.Bytes()))
	u2 := new(User)
	if err := BytesToInterface(byt, u2, "application/json"); err != nil {
		t.Error("Error", err.Error())
	}
	AssertEqual(u2.Age, 20, "Struct unmarshall", t)
	AssertEqual(u2.FirstName, "David", "Struct unmarshall", t)
	AssertEqual(u2.LastName, "Coperfield", "Struct unmarshall", t)

	//Now try users in array

	byt = bytes.NewBufferString(`[{"Id":"1","FirstName":"David","LastName":"Coperfield","Age":20,"Weight":0},{"Id":"2","FirstName":"David","LastName":"Coperfield","Age":20,"Weight":0}]`)
	userArr := make([]User, 0)

	if err := BytesToInterface(byt, &userArr, "application/json"); err != nil {
		t.Error("Error", err.Error())
	} else {
		for pos, au := range userArr {
			//println("User at pos:", pos," Data: ",au.FirstName,au.LastName,au.Id)
			AssertEqual(au.Id, strconv.Itoa(pos+1), "Slice unmarshall", t)
			AssertEqual(au.FirstName, "David", "Slice unmarshall", t)
		}
	}

	//Now try maps
	byt = bytes.NewBufferString(`{"One":{"Id":"One","FirstName":"Siya","LastName":"Dlamini","Age":29,"Weight":62},"Two":{"Id":"Two","FirstName":"Siya","LastName":"Dlamini","Age":29,"Weight":62}}`)
	userMap := make(map[string]User, 0)

	if err := BytesToInterface(byt, &userMap, "application/json"); err != nil {
		t.Error("Error", err.Error())
	} else {
		AssertEqual(userMap["One"].FirstName, "Siya", "Map Unmarshal", t)
		AssertEqual(userMap["One"].LastName, "Dlamini", "Map Unmarshal", t)

		AssertEqual(userMap["Two"].Id, "Two", "Map Unmarshal", t)
		AssertEqual(userMap["Two"].Age, 29, "Map Unmarshal", t)
	}

}
