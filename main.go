package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/atheerauribi/server1/proto"
	"google.golang.org/grpc"
)

// run proxy server, listen for requests from server1 and forward them to server2
func main() {
	go func() {
		if err := RunServer(); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	time.Sleep(4 * time.Second)
	// Define the server address
	serverAddr := flag.String("server_addr", "localhost:8888", "The server address in the format of host:port")

	fmt.Println("Dialling Client for Calculator service")
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewCalculatorClient(conn)

	// Call the Add method
	req := &pb.AddRequest{Number1: 21, Number2: 48}
	resp, err := client.Add(context.Background(), req)
	if err != nil {
		log.Fatalf("Add failed: %v", err)
	}
	fmt.Printf("Add: %v + %v = %v\n", req.Number1, req.Number2, resp.Result)
}
