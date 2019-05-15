package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"

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
	_, span := trace.StartSpan(ctx, "grpc: getping")
	span.Annotate([]trace.Attribute{
		trace.StringAttribute("message", in.Message),
	}, "todo:event")
	defer span.End()
	return &pb.Ping{Message: fmt.Sprintf("hi back! (%s)", in.Message)}, nil
}

//ProcessCommands processing a list of commands
func (s *Server) ProcessCommands(ctx context.Context, in *pb.CommandsRequest) (*pb.CuesResponse, error) {
	ctx, span := trace.StartSpan(ctx, "grpc: process comamnds")
	span.Annotate([]trace.Attribute{
		trace.Int64Attribute("num-commands", int64(len(in.Commands))),
	}, "todo:event")
	defer span.End()

	var responses []*pb.Cue
	m := s.master
	cs := m.GetDefaultCueStack()
	for _, eachCommand := range in.Commands {
		ctx, span := trace.StartSpan(ctx, "grpc: run single Command")
		x, err := cue.CommandToCue(ctx, m, eachCommand)
		if err != nil {
			tracing.SetError(span, err)
			span.End()
			return &pb.CuesResponse{Err: err.Error()}, err
		}
		x.Source.Input = cue.SourceInputRPC
		x.Source.Type = cue.SourceTypeCommand
		x.Source.Meta = eachCommand
		m.EnQueueCue(ctx, *x, cs)
		responses = append(responses, &pb.Cue{
			ExpectedDurationMS: int32(x.GetDuration() / time.Millisecond),
		})
		span.End()
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

		_, span := trace.StartSpan(stream.Context(), "grpc: streamGetLights")
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
		span.End()
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
		grpc.StatsHandler(new(ocgrpc.ServerHandler)),
	)
	pb.RegisterAPIServer(grpcServer, &Server{master: master})
	go grpcServer.Serve(lis)
	log.Printf("serving GRPC on %s", RPCConfig.Address)

	<-ctx.Done()
	log.Printf("[grpc] shutdown")
	grpcServer.GracefulStop()
	wg.Done()

}
