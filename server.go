package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/atheerauribi/server1/proto"
)

type userFacingServer struct {
	pb.UnimplementedCalculatorServer
	proxyClient pb.CalculatorClient
}

func (s *userFacingServer) connectToProxyServer() {
	// Create a client for the proxy grpc server at localhost:50052
	proxyAddr := "localhost:50052"
	conn, err := grpc.Dial(proxyAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	proxyClient := pb.NewCalculatorClient(conn)
	s.proxyClient = proxyClient
	fmt.Println("Connected to proxy server")
}

//Client calls proxyAdd, proxyAdd() packages Add request for calculator server, and returns the result

func (s *userFacingServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.OperationResponse, error) {
	// Call proxy server
	s.connectToProxyServer()
	outgoingReq := &pb.AddRequest{Number1: req.Number1, Number2: req.Number2}
	resp, err := s.proxyClient.Add(ctx, outgoingReq)
	if err != nil {
		log.Fatalf("could not call add in proxy: %v", err)
	}

	return resp, nil
}

// RunServer runs the gRPC server
func RunServer() error {
	fmt.Println("Running Server\nListening on port 8888...")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8888))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	s := grpc.NewServer(opts...)
	pb.RegisterCalculatorServer(s, &userFacingServer{})
	return s.Serve(lis)
}
