package peer

import (
	"context"
	"fmt"
	"log"
	pb "mutex/stc"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type Node struct {
	ID                string
	Address           string
	Peers             map[string]pb.MutexServiceClient
	RequestTimestamp  int64
	WantCS            bool
	InCS              bool
	DeferredResponses []string
	ResponseCount     int
	mu                sync.Mutex
	pb.UnimplementedMutexServiceServer
}

func NewNode(id, address string) *Node {
	return &Node{
		ID:                id,
		Address:           address,
		Peers:             make(map[string]pb.MutexServiceClient),
		DeferredResponses: make([]string, 0),
	}
}

func (n *Node) ConnectToPeer(peerID, peerAddr string) error {
	conn, err := grpc.NewClient(peerAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to peer %s: %v", peerID, err)
	}

	n.Peers[peerID] = pb.NewMutexServiceClient(conn)
	log.Printf("Node %s connected to peer %s at %s", n.ID, peerID, peerAddr)
	return nil
}

func (n *Node) RequestAccess(ctx context.Context, req *pb.AccessRequest) (*pb.AccessResponse, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	log.Printf("Node %s received request from %s with timestamp %d", n.ID, req.NodeId, req.Timestamp)

	if n.InCS || (n.WantCS && n.RequestTimestamp < req.Timestamp) {
		// Defer response
		n.DeferredResponses = append(n.DeferredResponses, req.NodeId)
		log.Printf("Node %s deferred response to %s", n.ID, req.NodeId)
		return &pb.AccessResponse{Granted: false}, nil
	}

	return &pb.AccessResponse{Granted: true}, nil
}

func (n *Node) ReleaseAccess(ctx context.Context, req *pb.ReleaseRequest) (*pb.ReleaseResponse, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	log.Printf("Node %s received release from %s", n.ID, req.NodeId)
	return &pb.ReleaseResponse{Acknowledged: true}, nil
}

func (n *Node) RequestCriticalSection() {
	n.mu.Lock()
	n.WantCS = true
	n.RequestTimestamp = time.Now().UnixNano()
	n.ResponseCount = 0
	n.mu.Unlock()

	log.Printf("Node %s requesting critical section access", n.ID)

	// Request access from all peers
	for peerID, client := range n.Peers {
		go func(id string, c pb.MutexServiceClient) {
			resp, err := c.RequestAccess(context.Background(), &pb.AccessRequest{
				NodeId:    n.ID,
				Timestamp: n.RequestTimestamp,
			})

			if err != nil {
				log.Printf("Error requesting access from %s: %v", id, err)
				return
			}

			if resp.Granted {
				n.mu.Lock()
				n.ResponseCount++
				if n.ResponseCount == len(n.Peers) {
					n.InCS = true
					n.ExecuteCriticalSection()
				}
				n.mu.Unlock()
			}
		}(peerID, client)
	}
}

func (n *Node) ExecuteCriticalSection() {
	log.Printf("Node %s entering critical section", n.ID)
	// Simulate critical section work
	time.Sleep(2 * time.Second)
	log.Printf("Node %s leaving critical section", n.ID)

	n.mu.Lock()
	n.InCS = false
	n.WantCS = false

	// Send deferred responses
	for _, peerID := range n.DeferredResponses {
		if client, ok := n.Peers[peerID]; ok {
			_, err := client.ReleaseAccess(context.Background(), &pb.ReleaseRequest{NodeId: n.ID})
			if err != nil {
				log.Printf("Error sending release to %s: %v", peerID, err)
			}
		}
	}

	n.DeferredResponses = make([]string, 0)
	n.mu.Unlock()
}
