package r

import (
	"fmt"
	"time"
)

func printTime() {
	fmt.Print(getTime(), " ")
}

var logTime = "02-01-2006 15:04:05.000"

func getTime() string {
	return time.Now().Format(logTime)
}

func log(a ...interface{}) {
	printTime()
	fmt.Println(a)
}
