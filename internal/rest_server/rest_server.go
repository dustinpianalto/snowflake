package rest_server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dustinpianalto/snowflake/internal/generator"
)

const (
	REST_PORT = ":50052"
)

type snowflake struct {
	Id    uint64 `json:"id"`
	IdStr string `json:"id_str"`
}

func getSnowflake(w http.ResponseWriter, r *http.Request) {
	outputChan := make(chan uint64, 1)
	defer close(outputChan)

	generator.Generator.RequestChan <- outputChan
	id := <-outputChan
	s := snowflake{
		Id:    id,
		IdStr: strconv.FormatUint(id, 10),
	}
	json.NewEncoder(w).Encode(s)
}

func RunRESTServer() {
	http.HandleFunc("/snowflake", getSnowflake)

	log.Printf("REST Server Listening on 0.0.0.0%s", REST_PORT)
	log.Fatal(http.ListenAndServe(REST_PORT, nil))
}
