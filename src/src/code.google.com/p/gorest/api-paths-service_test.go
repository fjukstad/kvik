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

type PathsService struct {
	RestService `root:"/paths-service/" consumes:"application/json" produces:"application/json" realm:"testing"`

	//Test Mixed paths with same length
	deleteMixed1 EndPoint `method:"DELETE" path:"/bool/{Bool:bool}/mix1/{Int:int}"`
	deleteMixed2 EndPoint `method:"DELETE" path:"/bool/{Bool:bool}/mix2/{Int:int}"`
	//Now check same path for different methods
	optionsMixed EndPoint `method:"OPTIONS" path:"/bool/{Bool:bool}/mix1/{Int:int}"`
	getMixed     EndPoint `method:"GET" path:"/bool/{Bool:bool}/mix1/{Int:int}" output:"string"`
	//getMixed2    EndPoint `method:"GET" path:"/bool/{Bool:bool}/mix1/{Int:int}" output:"string"`
}

func (serv PathsService) DeleteMixed1(Bool bool, Int int) {
	//Will return default response code of 200
}
func (serv PathsService) DeleteMixed2(Bool bool, Int int) {
	//Will return default response code of 200
}
func (serv PathsService) OptionsMixed(Bool bool, Int int) {
	rb := serv.ResponseBuilder()
	rb.Allow("GET")
	rb.Allow("HEAD").Allow("POST")
}
func (serv PathsService) GetMixed(Bool bool, Int int) string {
	return "Hello"
}
func (serv PathsService) GetMixed2(Bool bool, Int int) string {
	return "Hello"
}
