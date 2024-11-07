package main

import (
	"flag"
	"log"
	peer "mutex/peer"
	pb "mutex/stc"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
)

func main() {
	var (
		nodeID = flag.String("id", "", "Node ID")
		addr   = flag.String("addr", "", "Node address (host:port)")
		peers  = flag.String("peers", "", "Comma-separated list of peer addresses (id@host:port)")
	)
	flag.Parse()

	if *nodeID == "" || *addr == "" {
		log.Fatal("Node ID and address are required")
	}

	n := peer.NewNode(*nodeID, *addr)

	// Start gRPC server
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMutexServiceServer(grpcServer, n)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Connect to peers
	if *peers != "" {
		for _, peer := range strings.Split(*peers, ",") {
			parts := strings.Split(peer, "@")
			if len(parts) != 2 {
				log.Fatalf("Invalid peer format: %s", peer)
			}
			peerID, peerAddr := parts[0], parts[1]

			// Wait a bit for other nodes to start
			time.Sleep(time.Second * 2)
			if err := n.ConnectToPeer(peerID, peerAddr); err != nil {
				log.Printf("Failed to connect to peer %s: %v", peerID, err)
			}
		}
	}

	// Periodically request critical section access
	for {
		time.Sleep(1 * time.Second)
		n.RequestCriticalSection()
	}
}
