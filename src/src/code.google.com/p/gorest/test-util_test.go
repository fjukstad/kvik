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
	"testing"
)

func TestingAuthorizer(id string, role string) (bool, bool, SessionData) {
	if idsInRealm == nil {
		idsInRealm = make(map[string][]string, 0)
		idsInRealm["12345"] = []string{"var-user", "string-user", "post-user"}
		idsInRealm["fox"] = []string{"postInt-user"}
	}

	if roles, found := idsInRealm[id]; found {
		for _, r := range roles {
			if role == r {
				return true, true, nil
			}
		}
		return true, false, nil
	}

	return false, false, nil
}

func AssertEqual(given interface{}, expecting interface{}, compared string, t *testing.T) {
	if expecting != given {
		t.Error("Fail Assert:", compared, " Expecting:", expecting, "; but is:", given)
	} else {
		log.Println("Pass Assert:", compared)
	}
}
