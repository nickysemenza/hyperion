package client

import (
	"context"
	"log"

	pb "github.com/nickysemenza/hyperion/api/proto"
	"google.golang.org/grpc"
)

//Run runs the client
func Run() {
	conn, err := grpc.Dial("localhost:8888", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	client := pb.NewAPIClient(conn)

	res, err := client.GetPing(context.Background(), &pb.Ping{Message: "hi from client"})
	log.Println(res, err)

}
