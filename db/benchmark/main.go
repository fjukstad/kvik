package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fjukstad/kvik/db"
)

func main() {

	file, err := os.Create("cache.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	service := db.Service{"http://ec2-54-93-122-11.eu-central-1.compute.amazonaws.com:80", true}
	for range make([]int, 100) {
		start := time.Now()
		_, err := service.Summary("BRCA1")
		if err != nil {
			fmt.Println(err)
		}
		end := time.Now()
		duration := end.Sub(start)
		fmt.Fprintf(file, "%v\n", duration.Seconds())
	}
}
