package main

import (
	"context"
	"log"
	"time"

	pb "ninja/blog/rpc/blog"

	"google.golang.org/grpc"
)

const (
	Addr = "localhost:1235"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(Addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewArticleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Hello(ctx, &pb.HelloReq{Name: "wjh"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Resp: %s", r.Msg)
}
