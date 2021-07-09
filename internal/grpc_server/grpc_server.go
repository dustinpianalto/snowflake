package grpc_server

import (
	"context"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"

	"github.com/dustinpianalto/snowflake/internal/generator"
	"github.com/dustinpianalto/snowflake/snowflake"
)

const (
	GRPC_PORT = ":50051"
)

type SnowflakeServer struct {
	snowflake.UnimplementedSnowflakeServer
}

func (s *SnowflakeServer) GetSnowflake(ctx context.Context, in *snowflake.Empty) (*snowflake.SnowflakeReply, error) {
	var id uint64
	outputChan := make(chan uint64, 1)
	defer close(outputChan)
	generator.Generator.RequestChan <- outputChan
	id = <-outputChan
	return &snowflake.SnowflakeReply{Id: id, IdStr: strconv.FormatUint(id, 10)}, nil
}

func RunGRPCServer() {
	grpcListener, err := net.Listen("tcp", GRPC_PORT)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	snowflake.RegisterSnowflakeServer(grpcServer, &SnowflakeServer{})
	log.Printf("GRPC Server Listening on %v", grpcListener.Addr())
	if err := grpcServer.Serve(grpcListener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
