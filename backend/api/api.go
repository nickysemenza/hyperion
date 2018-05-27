package api

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/nickysemenza/hyperion/backend/proto"
	"google.golang.org/grpc"
)

//Server conforms to interface for proto generated stubs
type Server struct{}

//GetPing is test thing
func (s *Server) GetPing(ctx context.Context, in *pb.Ping) (*pb.Ping, error) {
	return &pb.Ping{Message: fmt.Sprintf("hi back! (%s)", in.Message)}, nil
}

//ServerRPC runs a RPC server
func ServeRPC(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAPIServer(grpcServer, &Server{})
	grpcServer.Serve(lis)
}
