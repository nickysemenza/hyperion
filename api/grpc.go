package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"

	pb "github.com/nickysemenza/hyperion/api/proto"
	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/cue"
	"github.com/nickysemenza/hyperion/core/light"
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

func (s *Server) StreamGetLights(in *pb.Empty, stream pb.API_StreamGetLightsServer) error {
	for {
		allLights := light.GetLightsByName()
		var pbLights []*pb.Light

		for k, v := range allLights {
			color := light.GetCurrentState(v.GetName()).RGB.AsPB()
			pbLights = append(pbLights, &pb.Light{
				Name:         k,
				Type:         v.GetType(),
				CurrentColor: &color,
			})
		}

		err := stream.Send(&pb.Lights{
			Lights: pbLights,
		})
		if err != nil {
			log.Println(err)
			break
		}
		time.Sleep(time.Second)
	}
	return nil
}

//ServeRPC runs a RPC server
func ServeRPC(ctx context.Context) {
	RPCConfig := config.GetServerConfig(ctx).Inputs.RPC
	if !RPCConfig.Enabled {
		log.Info("rpc is not enabled")
		return
	}
	lis, err := net.Listen("tcp", RPCConfig.Address)
	// lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAPIServer(grpcServer, &Server{})
	grpcServer.Serve(lis)
}
