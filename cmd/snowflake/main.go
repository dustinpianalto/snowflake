package main

import (
	"sync"

	"github.com/dustinpianalto/snowflake/internal/generator"
	"github.com/dustinpianalto/snowflake/internal/grpc_server"
	"github.com/dustinpianalto/snowflake/internal/rest_server"
)

func main() {

	generator.CreateGenerator(1)

	go generator.Generator.Run()

	var wg sync.WaitGroup
	wg.Add(1)
	go grpc_server.RunGRPCServer()
	go rest_server.RunRESTServer()
	wg.Wait()
}
