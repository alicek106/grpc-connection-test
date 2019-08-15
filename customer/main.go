package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	pb "github.com/alicek106/grpc-connection-test/messages"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	addr := ":8080"
	orderAddr := os.Getenv("ORDER_ADDR") // Order service DNS
	conn, err := grpc.Dial(orderAddr, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("gRPC server dial failed: %v", err)
	}

	client := pb.NewOrderingClient(conn)

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		stuff := query.Get("stuff")
		i, err := strconv.ParseInt(query.Get("money"), 10, 32)
		if err != nil {
			log.Fatalf("Failed to parse int value", err)
		}

		money := int32(i)

		ctx := req.Context()
		resp, err := client.Order(ctx, &pb.OrderRequest{
			Stuff: stuff,
			Money: money,
		})
		if err != nil {
			log.Fatalf("Failed to call RPC Order", err)
		}

		log.Infof("Stuff: %s, Change: %d", resp.GetStuff(), resp.GetChange())
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(fmt.Sprintf("Served from %s", resp.GetIp())))
	})

	http.HandleFunc("/healthz", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	log.Infof("HTTP server listening on %s", addr)
	http.ListenAndServe(addr, nil)
}
