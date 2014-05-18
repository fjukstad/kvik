package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fjukstad/rpcman"
)

func genrands(num int) (ret []float64) {

	ret = make([]float64, num)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		ret[i] = r.Float64() * 100
	}
	return ret

}

func main() {
	rpc, err := rpcman.Init("tcp://localhost:5555")
	if err != nil {
		return
	}

	for i := 0; i < 10; i++ {
		floats := genrands(100)
		sum, _ := rpc.Call("sum", floats)
		std, _ := rpc.Call("std", floats)
		variance, _ := rpc.Call("var", floats)
		mean, _ := rpc.Call("mean", floats)

		add, _ := rpc.Call("add", 2, 3)

		fmt.Println(sum, std, variance, mean, add)
	}

	// rpc.Close()

}
