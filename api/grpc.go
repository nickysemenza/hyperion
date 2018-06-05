package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/nickysemenza/hyperion/api/proto"
	"github.com/nickysemenza/hyperion/cue"
	"google.golang.org/grpc"
)

//Server conforms to interface for proto generated stubs
type Server struct{}

//GetPing is test thing
func (s *Server) GetPing(ctx context.Context, in *pb.Ping) (*pb.Ping, error) {
	return &pb.Ping{Message: fmt.Sprintf("hi back! (%s)", in.Message)}, nil
}

func (s *Server) StreamCueMaster(in *pb.Ping, stream pb.API_StreamCueMasterServer) error {
	log.Println("StreamCueMaster started")
	for {
		cm := cue.GetCueMaster()
		bytes, _ := json.Marshal(cm)

		err := stream.Send(&pb.MarshalledJSON{Data: bytes})
		if err != nil {
			log.Println(err)
			break
		}
		time.Sleep(time.Second)
	}
	return nil

}

//ServeRPC runs a RPC server
func ServeRPC(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAPIServer(grpcServer, &Server{})
	grpcServer.Serve(lis)
}
