package main

import "fmt"

func main() {
	var (
		err error
	)
	if err = startES(); err != nil {
		panic(err)
	}
	fmt.Println(esClient)
}
