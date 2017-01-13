package main

import (
	"fmt"

	"github.com/fjukstad/kvik/r"
)

func main() {

	session, err := r.NewSession(-1)
	if err != nil {
		fmt.Println(err)
		return
	}
	key, err := session.Call("stats", "rnorm", "n=100")
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := session.Get(key, "json")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(res))

	res, err = session.Rpc("stats", "rnorm", "n=200", "json")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(res))

}
