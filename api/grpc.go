package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
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

//StreamCueMaster streams the cuemaster
func (s *Server) StreamCueMaster(in *pb.ConnectionSettings, stream pb.API_StreamCueMasterServer) error {
	tick, err := time.ParseDuration(in.Tick)
	if err != nil {
		return err
	}
	log.Println("StreamCueMaster started")
	for {
		cm := cue.GetCueMaster()
		bytes, _ := json.Marshal(cm)

		err := stream.Send(&pb.MarshalledJSON{Data: bytes})
		if err != nil {
			log.Println(err)
			break
		}
		time.Sleep(tick)
	}
	return nil

}

//StreamGetLights sends the light state to the client on an interval
func (s *Server) StreamGetLights(in *pb.ConnectionSettings, stream pb.API_StreamGetLightsServer) error {
	tick, err := time.ParseDuration(in.Tick)
	if err != nil {
		return err
	}
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
		time.Sleep(time.Millisecond * 20)
	}
	return nil
}

//ServeRPC runs a RPC server
func ServeRPC(ctx context.Context, wg *sync.WaitGroup) {
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
	go grpcServer.Serve(lis)

	<-ctx.Done()
	log.Printf("[grpc] shutdown")
	grpcServer.GracefulStop()
	wg.Done()

}
