package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"sort"

	"github.com/aybabtme/rgbterm"
	pb "github.com/nickysemenza/hyperion/api/proto"
	"google.golang.org/grpc"
)

//Run runs the client
func Run(address string) {
	conn, cerr := grpc.Dial(address, grpc.WithInsecure())
	if cerr != nil {
		log.Println(cerr)
	}
	defer conn.Close()

	client := pb.NewAPIClient(conn)

	lights := make(map[string]*pb.Light)
	stream, err := client.StreamGetLights(context.Background(), &pb.Empty{})
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
