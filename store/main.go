package main

import (
	"context"
	"net"
	"os"

	pb "github.com/alicek106/grpc-connection-test/messages"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	addr := ":80" // os.Getenv("ADDR")

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to initialize TCP listen: %v", err)
	}
	defer lis.Close()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed to read hostname: %v", err)
	}
	storeName := storeServer{storeID: hostname}

	server := grpc.NewServer()
	pb.RegisterOrderingServer(server, storeName)

	log.Infof("Store initialized, storeID = %s", storeName.storeID)
	log.Printf("gRPC Listening on %s", lis.Addr().String())

	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}

type storeServer struct {
	storeID string
}

func (ss storeServer) Order(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
	stuff := req.GetStuff()
	money := req.GetMoney()

	return &pb.OrderResponse{
		Stuff:  stuff,
		Ip:     ss.storeID,
		Change: int32(money),
	}, nil
}
