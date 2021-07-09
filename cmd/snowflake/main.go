package main

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/dustinpianalto/snowflake/internal/generator"
	"github.com/dustinpianalto/snowflake/internal/grpc_server"
	"github.com/dustinpianalto/snowflake/internal/rest_server"
)

func main() {

	workerID, err := strconv.ParseUint(os.Getenv("WORKER_ID"), 10, 64)
	if err != nil {
		log.Fatal("Not a valid worker id")
	}

	generator.CreateGenerator(workerID)

	go generator.Generator.Run()

	var wg sync.WaitGroup
	wg.Add(1)
	go grpc_server.RunGRPCServer()
	go rest_server.RunRESTServer()
	wg.Wait()
}
