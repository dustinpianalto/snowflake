package main

import (
	"fmt"

	"github.com/dustinpianalto/snowflake/internal/generator"
)

func main() {

	gen, err := generator.CreateGenerator(1)
	if err != nil {
		panic(err)
	}

	requestChan := make(chan chan uint64)
	outputChan := make(chan uint64)

	go func() {
		for id := range outputChan {
			fmt.Println(id)
		}
	}()
	go gen.Run(requestChan)

	for i := 0; i < 500000; i++ {
		requestChan <- outputChan
	}
}
