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
	"reflect"
	"strings"
)

type argumentData struct {
	parameter param
	data      string
}
type param struct {
	positionInPath int
	name           string
	typeName       string
}

var aLLOWED_PAR_TYPES = []string{"string", "int", "int32", "int64", "bool", "float32", "float64"}

func prepServiceMetaData(root string, tags reflect.StructTag, i interface{}, name string) serviceMetaData {
	md := new(serviceMetaData)

	if tag := tags.Get("root"); tag != "" {
		md.root = tag
	}
	if root != "" {
		md.root = root + md.root
	}
	log.Println("All EndPoints for service [", name, "] , registered under root path: ", md.root)
	if tag := tags.Get("consumes"); tag != "" {
		md.consumesMime = tag
		if GetMarshallerByMime(tag) == nil {
			log.Panic("The Marshaller for mime-type:[" + tag + "], is not registered. Please register this type before registering your service.")
		}

	} else {
		md.consumesMime = Application_Json //Default
	}
	if tag := tags.Get("produces"); tag != "" {
		md.producesMime = tag
		if GetMarshallerByMime(tag) == nil {
			log.Panic("The Marshaller for mime-type:[" + tag + "], is not registered. Please register this type before registering your service.")
		}
	} else {
		md.producesMime = Application_Json //Default
	}

	if tag := tags.Get("realm"); tag != "" {
		md.realm = tag
		if GetAuthorizer(tag) == nil {
			log.Panic("The realm:[" + tag + "], is not registered. Please register this realm before registering your service.")
		}
	}

	md.template = i
	return *md
}

func makeEndPointStruct(tags reflect.StructTag, serviceRoot string) endPointStruct {

	ms := new(endPointStruct)

	if tag := tags.Get("method"); tag != "" {
		if tag == "GET" {
			ms.requestMethod = GET
		} else if tag == "POST" {
			ms.requestMethod = POST
		} else if tag == "PUT" {
			ms.requestMethod = PUT
		} else if tag == "DELETE" {
			ms.requestMethod = DELETE
		} else if tag == "HEAD" {
			ms.requestMethod = HEAD
		} else if tag == "OPTIONS" {
			ms.requestMethod = OPTIONS
		} else {
			log.Panic("Unknown method type:[" + tag + "] in endpoint declaration. Allowed types {GET,POST,PUT,DELETE,HEAD,OPTIONS}")
		}

		if tag := tags.Get("path"); tag != "" {
			serviceRoot = strings.TrimRight(serviceRoot, "/")
			ms.signiture = serviceRoot + "/" + strings.Trim(tag, "/")
		} else {
			log.Panic("Endpoint declaration must have the tags 'method' and 'path' ")
		}

		if tag := tags.Get("output"); tag != "" {
			ms.outputType = tag
			if strings.HasPrefix(tag, "[]") { //Check for slice/array/list types.
				ms.outputTypeIsArray = true
				ms.outputType = ms.outputType[2:]
			}
			if strings.HasPrefix(tag, "map[") { //Check for map[string]. We only handle string keyed maps!!!

				if ms.outputType[4:10] == "string" {
					ms.outputTypeIsMap = true
					ms.outputType = ms.outputType[11:]
				} else {
					log.Panic("Only string keyed maps e.g( map[string]... ) are allowed on the [output] tag. Endpoint: " + ms.signiture)
				}

			}
		}

		if tag := tags.Get("input"); tag != "" {
			ms.inputMime = tag
		}
		if tag := tags.Get("role"); tag != "" {
			ms.role = tag
		}

		if tag := tags.Get("postdata"); tag != "" {
			ms.postdataType = tag
			if strings.HasPrefix(tag, "[]") { //Check for slice/array/list types.
				ms.postdataTypeIsArray = true
				ms.postdataType = ms.postdataType[2:]
			}
			if strings.HasPrefix(tag, "map[") { //Check for map[string]. We only handle string keyed maps!!!

				if ms.postdataType[4:10] == "string" {
					ms.postdataTypeIsMap = true
					ms.postdataType = ms.postdataType[11:]
				} else {
					log.Panic("Only string keyed maps e.g( map[string]... ) are allowed on the [postdata] tag. Endpoint: " + ms.signiture)
				}

			}
		}

		parseParams(ms)
		return *ms
	}

	log.Panic("Endpoint declaration must have the tags 'method' and 'path' ")
	return *ms //Should not get here

}
func parseParams(e *endPointStruct) {
	e.signiture = strings.Trim(e.signiture, "/")
	e.params = make([]param, 0)
	e.queryParams = make([]param, 0)
	e.nonParamPathPart = make(map[int]string, 0)

	pathPart := e.signiture
	queryPart := ""

	if i := strings.Index(e.signiture, "?"); i != -1 {

		pathPart = e.signiture[:i]
		//e.root = pathPart
		pathPart = strings.TrimRight(pathPart, "/")
		queryPart = e.signiture[i+1:]

		//Extract Query Parameters

		for pos, str1 := range strings.Split(queryPart, "&") {
			if strings.HasPrefix(str1, "{") && strings.HasSuffix(str1, "}") {
				parName, typeName := getVarTypePair(str1, e.signiture)

				for _, par := range e.queryParams {
					if par.name == parName {
						log.Panic("Duplicate Query Parameter name(" + parName + ") in REST path: " + e.signiture)
					}
				}
				//e.queryParams[len(e.queryParams)] = param{pos, parName, typeName}
				e.queryParams = append(e.queryParams, param{pos, parName, typeName})
			} else {
				log.Panic("Please check that your Query Parameters are configured correctly for endpoint: " + e.signiture)
			}
		}
	}

	if i := strings.Index(pathPart, "{"); i != -1 {
		e.root = pathPart[:i]
	} else {
		e.root = pathPart
	}

	//Extract Path Parameters
	for pos, str1 := range strings.Split(pathPart, "/") {
		e.signitureLen++

		if strings.HasPrefix(str1, "{") && strings.HasSuffix(str1, "}") { //This just ensures we re dealing with a varibale not normal path.

			parName, typeName := getVarTypePair(str1, e.signiture)

			if parName == "..." {
				e.isVariableLength = true
				parName, typeName := getVarTypePair(str1, e.signiture)
				e.params = append(e.params, param{pos, parName, typeName})
				e.paramLen++
				break
			}
			for _, par := range e.params {
				if par.name == parName {
					log.Panic("Duplicate Path Parameter name(" + parName + ") in REST path: " + e.signiture)
				}
			}

			e.params = append(e.params, param{pos, parName, typeName})
			e.paramLen++
		} else {
			e.nonParamPathPart[pos] = str1

		}
	}

	e.root = strings.TrimRight(e.root, "/")

	if e.isVariableLength && e.paramLen > 1 {
		log.Panic("Variable length endpoints can only have one parameter declaration: " + pathPart)
	}

	for key, ep := range _manager().endpoints {
		if ep.root == e.root && ep.signitureLen == e.signitureLen && reflect.DeepEqual(ep.nonParamPathPart, e.nonParamPathPart) && ep.requestMethod == e.requestMethod {
			log.Panic("Can not register two endpoints with same request-method(" + ep.requestMethod + ") and same signature: " + e.signiture + " VS " + ep.signiture)
		}
		if ep.requestMethod == e.requestMethod && pathPart == key {
			log.Panic("Endpoint already registered: " + pathPart)
		}
		if e.isVariableLength && (strings.Index(ep.root+"/", e.root+"/") == 0 || strings.Index(e.root+"/", ep.root+"/") == 0) && ep.requestMethod == e.requestMethod {
			log.Panic("Variable length endpoints can only be mounted on a unique root. Root already used: " + ep.root + " <> " + e.root)
		}
	}
}

func getVarTypePair(part string, sign string) (parName string, typeName string) {

	temp := strings.Trim(part, "{}")
	ind := 0
	if ind = strings.Index(temp, ":"); ind == -1 {
		log.Panic("Please ensure that parameter names(" + temp + ") have associated types in REST path: " + sign)
	}
	parName = temp[:ind]
	typeName = temp[ind+1:]

	if !isAllowedParamType(typeName) {
		log.Panic("Type " + typeName + " is not allowed for Path/Query-parameters in REST path: " + sign)
	}

	return
}

func isAllowedParamType(typeName string) bool {
	for _, s := range aLLOWED_PAR_TYPES {
		if s == strings.ToLower(typeName) {
			return true
		}
	}
	return false
}

func getEndPointByUrl(method string, url string) (endPointStruct, map[string]string, map[string]string, string, bool) {
	//println("Getting:",url)

	pathPart := url
	queryPart := ""

	if i := strings.Index(url, "?"); i != -1 {
		pathPart = url[:i]
		queryPart = url[i+1:]
	}

	pathPart = strings.Trim(pathPart, "/")
	totalParts := strings.Count(pathPart, "/")
	totalParts++

	epRet := new(endPointStruct)
	pathArgs := make(map[string]string, 0)
	queryArgs := make(map[string]string, 0)

	var ep *endPointStruct

EPLOOP:
	for _, loopEp := range _manager().endpoints {
		//              println(method, ":", loopEp.requestMethod, pathPart, ":", loopEp.root, totalParts, ":", loopEp.signitureLen, "Variable?", loopEp.isVariableLength)
		if loopEp.isVariableLength && (strings.Index(pathPart+"/", loopEp.root+"/") == 0) && loopEp.requestMethod == method {
			ep = &loopEp
			varsPart := strings.Trim(pathPart[len(loopEp.root):], "/")
			//                      println("::::::::::::::::Root", pathPart, ">>>>>>>Vars", varsPart)
			for upos, str1 := range strings.Split(varsPart, "/") {
				pathArgs[string(upos)] = strings.Trim(str1, " ")
			}
		} else if (strings.Index(pathPart+"/", loopEp.root+"/") == 0) && loopEp.signitureLen == totalParts && loopEp.requestMethod == method {
			ep = &loopEp
			//We first make sure that the other parts of the path that are not parameters do actully match with the signature.
			//If not we exit. We do not have to cary on looking since we only allow one registration per root and length.
			for pos, name := range ep.nonParamPathPart {
				for upos, str1 := range strings.Split(pathPart, "/") {
					if upos == pos {
						if name != str1 {
							//Even though the beginning of the path matched, some other part didn't, keep looking.
							ep = nil
							continue EPLOOP
						}
						break
					}
				}
			}
			//Extract Path Arguments
			for _, par := range ep.params {
				for upos, str1 := range strings.Split(pathPart, "/") {

					if par.positionInPath == upos {
						pathArgs[par.name] = strings.Trim(str1, " ")
						break
					}
				}
			}
		}

		if ep != nil {
			xsrft := ""
			//Extract Query Arguments: These are optional in the query, so some or all of them might not be there.
			//Also, if they are there, they do not have to be in the same order they were sepcified in on the declaration signature.
			for _, str1 := range strings.Split(queryPart, "&") {
				if i := strings.Index(str1, "="); i != -1 {
					pName := str1[:i]
					dataString := str1[i+1:]
					if pName == XSXRF_PARAM_NAME {
						xsrft = strings.Trim(dataString, " ")
						log.Print("Session Id:", xsrft)
					} else {
						for _, par := range ep.queryParams {
							if par.name == pName {
								queryArgs[pName] = strings.Trim(dataString, " ")
								break
							}
						}
					}

				}
			}

			return *ep, pathArgs, queryArgs, xsrft, true //Path found
		}
	}

	return *epRet, pathArgs, queryArgs, "", false //Path not found
}
