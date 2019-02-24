package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	pb "github.com/nickysemenza/hyperion/api/proto"
	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/cue"
	"github.com/nickysemenza/hyperion/util/tracing"
	"google.golang.org/grpc"
)

//Server conforms to interface for proto generated stubs
type Server struct {
	master cue.MasterManager
}

//GetPing is test thing
func (s *Server) GetPing(ctx context.Context, in *pb.Ping) (*pb.Ping, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "grpc: getping")
	span.LogKV("message", in.Message)
	defer span.Finish()
	return &pb.Ping{Message: fmt.Sprintf("hi back! (%s)", in.Message)}, nil
}

//ProcessCommands processing a list of commands
func (s *Server) ProcessCommands(ctx context.Context, in *pb.CommandsRequest) (*pb.CuesResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "grpc: process comamnds")
	span.LogKV("num-commands", len(in.Commands))
	defer span.Finish()

	var responses []*pb.Cue
	m := s.master
	cs := m.GetDefaultCueStack()
	for _, eachCommand := range in.Commands {
		span, ctx := opentracing.StartSpanFromContext(ctx, "grpc: run single Command")
		x, err := cue.CommandToCue(ctx, m, eachCommand)
		if err != nil {
			tracing.SetError(span, err)
			span.Finish()
			return &pb.CuesResponse{Err: err.Error()}, err
		}
		x.Source.Input = cue.SourceInputRPC
		x.Source.Type = cue.SourceTypeCommand
		x.Source.Meta = eachCommand
		m.EnQueueCue(ctx, *x, cs)
		responses = append(responses, &pb.Cue{
			ExpectedDurationMS: int32(x.GetDuration() / time.Millisecond),
		})
		span.Finish()
	}
	return &pb.CuesResponse{Cues: responses}, nil
}

//StreamCueMaster streams the cuemaster
func (s *Server) StreamCueMaster(in *pb.ConnectionSettings, stream pb.API_StreamCueMasterServer) error {
	tick, err := time.ParseDuration(in.Tick)
	if err != nil {
		return err
	}
	log.Println("StreamCueMaster started")
	for {
		bytes, _ := json.Marshal(s.master)

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

		span, _ := opentracing.StartSpanFromContext(stream.Context(), "grpc: streamGetLights")
		allLights := s.master.GetLightManager().GetLightsByName() //TODO: fix this
		var pbLights []*pb.Light

		for k, v := range allLights {
			color := s.master.GetLightManager().GetState(v.GetName()).RGB.AsPB()
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
		span.Finish()
		time.Sleep(tick)
	}
	return nil
}

//ServeRPC runs a RPC server
func ServeRPC(ctx context.Context, wg *sync.WaitGroup, master cue.MasterManager) {
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
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_opentracing.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
		)),
	)
	pb.RegisterAPIServer(grpcServer, &Server{master: master})
	go grpcServer.Serve(lis)

	<-ctx.Done()
	log.Printf("[grpc] shutdown")
	grpcServer.GracefulStop()
	wg.Done()

}
