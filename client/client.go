package client

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/nickysemenza/hyperion/core/cue"

	pb "github.com/nickysemenza/hyperion/api/proto"
	"google.golang.org/grpc"
)

//Run runs the client
func Run(address string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	client := pb.NewAPIClient(conn)

	go func() {
		time.Sleep(2 * time.Second)
		res, err := client.GetPing(context.Background(), &pb.Ping{Message: "hi from client"})
		log.Println(res, err)
	}()

	stream, err := client.StreamCueMaster(context.Background(), &pb.Ping{Message: "hi from client"})
	for {
		received, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(client, err)
		}

		var cm cue.Master

		json.Unmarshal(received.Data, &cm)

		spew.Dump(cm)
	}

}
