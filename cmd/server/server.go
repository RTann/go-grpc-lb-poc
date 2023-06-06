// server.go
package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pb "github.com/jvdm/go-grpc-lb-poc/generated/api"
	"github.com/jvdm/go-grpc-lb-poc/metrics"
)

type server struct {
	pb.UnimplementedPocServiceServer
	serverID  string
	startTime int64
}

func NewServer(id string) *server {
	return &server{
		serverID:  id,
		startTime: time.Now().Unix(),
	}
}

func (s *server) SendRequest(ctx context.Context, req *pb.PocRequest) (*pb.PocResponse, error) {
	start := time.Now()
	log.Printf("Request: %s", start)
	defer func() {
		processingTime := time.Since(start).Seconds()
		metrics.ServerResponseTimes.WithLabelValues(s.serverID).Observe(processingTime)
		metrics.ServerRequestsReceived.WithLabelValues(s.serverID).Inc()
	}()

	// Simulate processing time
	time.Sleep(time.Second)

	return &pb.PocResponse{
		ServerId:        s.serverID,
		ProcessingTime:  start.Unix() - req.Timestamp,
		ServerStartTime: s.startTime,
	}, nil
}

func main() {
	serverID, err := os.Hostname()
	if err != nil {
		log.Fatalf("getting hostname: %v", err)
	}

	listen, err := net.Listen("tcp", ":5000")
	if err != nil {
		metrics.ServerErrors.WithLabelValues(serverID).Inc()
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionAge: 10 * time.Minute,
	}))
	pb.RegisterPocServiceServer(s, NewServer(serverID))

	// Start metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	log.Println("Starting server...")
	if err := s.Serve(listen); err != nil {
		metrics.ServerErrors.WithLabelValues(serverID).Inc()
		log.Fatalf("failed to serve: %v", err)
	}
}
