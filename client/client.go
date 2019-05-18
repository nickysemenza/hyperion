package client

import (
	"context"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"

	pb "github.com/nickysemenza/hyperion/api/proto"
	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/util/tracing"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//Run runs the client
func Run(ctx context.Context) {
	config := config.GetClientConfig(ctx)
	spew.Dump(config)
	tracing.InitTracer(config.Tracing.ServerAddress, "hyperion-client")

	conn, cerr := grpc.Dial(config.ServerAddress, grpc.WithInsecure(),
		grpc.WithStatsHandler(new(ocgrpc.ClientHandler)))
	if cerr != nil {
		log.Println(cerr)
	}
	defer conn.Close()

	client := pb.NewAPIClient(conn)

	sCtx, span := trace.StartSpan(ctx, "getping")
	ping, _ := client.GetPing(sCtx, &pb.Ping{Message: "test"})
	span.Annotate([]trace.Attribute{}, "error")
	span.End()

	sCtx, span = trace.StartSpan(ctx, "command test")
	resp, err := client.ProcessCommands(sCtx, &pb.CommandsRequest{Commands: []string{
		"set(par1+hue2:green+red:0.7s+1s)",
		"set(par1:blue:1s)",
		"blackout(0)",
		"set(par1:#FF5500:300ms)",
	}})
	if err != nil {
		tracing.SetError(span, err)
	}
	spew.Dump(resp)
	span.Annotate([]trace.Attribute{
		trace.StringAttribute("resp", ping.GetMessage()),
	}, "error")
	span.End()

	// lights := make(map[string]*pb.Light)
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
	// stream, err := client.StreamGetLights(ctx, &pb.ConnectionSettings{Tick: "20ms"})
	// if err != nil {
	// 	log.Fatal(client, err)
	// }
	// for {
	// 	received, err := stream.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatal(client, err)
	// 	}

	// 	// spew.Dump(received)
	// 	for _, l := range received.Lights {
	// 		lights[l.Name] = l
	// 	}

	// 	var keys []string
	// 	for k := range lights {
	// 		keys = append(keys, k)
	// 	}
	// 	sort.Strings(keys)

	// 	for _, k := range keys {
	// 		rgb := lights[k].CurrentColor
	// 		colorBlock := rgbterm.Bytes([]byte("███"), uint8(rgb.R), uint8(rgb.G), uint8(rgb.B), 0, 0, 0)
	// 		// fmt.Printf("%s, %v, %s\n", k, lights[k])
	// 		fmt.Printf("%s %s\n", colorBlock, k)
	// 	}
	// 	fmt.Println("-----")
	// }

}
