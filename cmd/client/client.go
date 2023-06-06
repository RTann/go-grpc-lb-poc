// client.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/jvdm/go-grpc-lb-poc/generated/api"
	"github.com/jvdm/go-grpc-lb-poc/metrics"
)

type serverInfo struct {
	lastResponseTime time.Time
	startTime        time.Time
}

func main() {
	serverHostname, ok := os.LookupEnv("SERVER_HOSTNAME")
	if !ok {
		serverHostname = "localhost"
	}
	log.Printf("Connecting to server: %s...", serverHostname)
	conn, err := grpc.Dial("dns:///"+serverHostname+":5000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"[{"round_robin": {}}]"}`),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Printf("Connected.")
	defer conn.Close()

	client := pb.NewPocServiceClient(conn)

	// Start metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	// Initialize known servers
	knownServers := make(map[string]*serverInfo)

	for {
		start := time.Now()
		log.Printf("Request: %s", start)
		response, err := client.SendRequest(context.Background(), &pb.PocRequest{
			ClientId:  "client-1",
			Timestamp: start.Unix(),
		})
		log.Printf("Response: %s", response)
		if err != nil {
			log.Printf("Could not send request: %v", err)
			metrics.ClientErrors.Inc()
		} else {
			latency := time.Since(start).Seconds()
			metrics.ClientRequestLatencies.WithLabelValues(response.ServerId).Observe(latency)
			metrics.ClientRequestDistribution.WithLabelValues(response.ServerId).Inc()
			metrics.ClientRequestsSent.WithLabelValues(response.ServerId).Inc()

			// Check for new servers
			info, ok := knownServers[response.ServerId]
			if !ok {
				info = &serverInfo{
					startTime: time.Unix(response.ServerStartTime, 0),
				}
				knownServers[response.ServerId] = info
				serverStartTime := time.Unix(response.ServerStartTime, 0)
				detectionTime := start.Sub(serverStartTime).Seconds()
				metrics.ClientChangeDetectionTimes.WithLabelValues(response.ServerId).Observe(detectionTime)
			}
			info.lastResponseTime = time.Now()
		}

		time.Sleep(time.Second)

		for serverID, info := range knownServers {
			inactivityTime := time.Since(info.lastResponseTime).Seconds()
			metrics.ClientInactivityTimes.WithLabelValues(serverID).Set(inactivityTime)
		}
	}
}
