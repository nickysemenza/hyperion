package client

import (
	"context"
	"fmt"
	"io"
	"sort"

	pb "github.com/nickysemenza/hyperion/api/proto"
	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/util/tracing"

	"github.com/aybabtme/rgbterm"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//Run runs the client
func Run(ctx context.Context) {
	config := config.GetClientConfig(ctx)

	tracing.InitTracer(config.Tracing.ServerAddress, "hyperion-client")

	conn, cerr := grpc.Dial(config.ServerAddress, grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_opentracing.StreamClientInterceptor(),
		)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(),
		)))
	if cerr != nil {
		log.Println(cerr)
	}
	defer conn.Close()

	client := pb.NewAPIClient(conn)

	span, _ := opentracing.StartSpanFromContext(ctx, "getping")
	ping, _ := client.GetPing(ctx, &pb.Ping{Message: "test"})
	span.LogKV("resp", ping.GetMessage())
	span.Finish()

	lights := make(map[string]*pb.Light)

	stream, err := client.StreamGetLights(ctx, &pb.ConnectionSettings{Tick: "20ms"})
	if err != nil {
		log.Fatal(client, err)
	}
	for {
		received, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(client, err)
		}

		// spew.Dump(received)
		for _, l := range received.Lights {
			lights[l.Name] = l
		}

		var keys []string
		for k := range lights {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			rgb := lights[k].CurrentColor
			colorBlock := rgbterm.Bytes([]byte("███"), uint8(rgb.R), uint8(rgb.G), uint8(rgb.B), 0, 0, 0)
			// fmt.Printf("%s, %v, %s\n", k, lights[k])
			fmt.Printf("%s %s\n", colorBlock, k)
		}
		fmt.Println("-----")
	}

}
