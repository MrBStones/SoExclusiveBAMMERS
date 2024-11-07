package peer

import (
	"context"
	"fmt"
	"log"
	pb "mutex/stc"
	"sort"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Request struct {
	NodeID    string
	Timestamp uint64
}

type Node struct {
	ID                string
	Address           string
	Peers             map[string]pb.MutexServiceClient
	LamportClock      uint64
	WantCS            bool
	InCS              bool
	RequestQueue      []Request
	DeferredResponses map[string]bool
	ResponseCount     int
	CurrentRequest    *Request
	mu                sync.Mutex
	pb.UnimplementedMutexServiceServer
}

func NewNode(id, address string) *Node {
	return &Node{
		ID:                id,
		Address:           address,
		Peers:             make(map[string]pb.MutexServiceClient),
		LamportClock:      0,
		RequestQueue:      make([]Request, 0),
		DeferredResponses: make(map[string]bool),
	}
}

func (n *Node) UpdateLamportClock(msgTimestamp uint64) uint64 {

	// Same as: clock = max(local_clock, msg_timestamp) + 1
	if msgTimestamp > n.LamportClock {
		n.LamportClock = msgTimestamp
	}
	n.LamportClock++
	return n.LamportClock
}

func (n *Node) GetLamportClock() uint64 {
	n.LamportClock++
	return n.LamportClock
}

func (n *Node) ConnectToPeer(peerID, peerAddr string) error {

	conn, err := grpc.NewClient(peerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to peer %s: %v", peerID, err)
	}

	n.Peers[peerID] = pb.NewMutexServiceClient(conn)
	log.Printf("Node %s connected to peer %s at %s", n.ID, peerID, peerAddr)
	return nil
}

func (n *Node) RequestAccess(ctx context.Context, req *pb.AccessRequest) (*pb.AccessResponse, error) {

	// Update Lamport clock on message receipt
	timestamp := n.UpdateLamportClock(req.LamportTimestamp)

	log.Printf("Node %s received request from %s with Lamport timestamp %d", n.ID, req.NodeId, req.LamportTimestamp)

	// Add request to queue
	n.RequestQueue = append(n.RequestQueue, Request{
		NodeID:    req.NodeId,
		Timestamp: req.LamportTimestamp,
	})

	// Sort queue by timestamp and node ID
	sort.Slice(n.RequestQueue, func(i, j int) bool {
		if n.RequestQueue[i].Timestamp == n.RequestQueue[j].Timestamp {
			return n.RequestQueue[i].NodeID < n.RequestQueue[j].NodeID
		}
		return n.RequestQueue[i].Timestamp < n.RequestQueue[j].Timestamp
	})

	if n.InCS || (n.WantCS && n.isHigherPriority(n.CurrentRequest, &Request{
		NodeID:    req.NodeId,
		Timestamp: req.LamportTimestamp,
	})) {
		n.DeferredResponses[req.NodeId] = true
		log.Printf("Node %s deferring response to %s", n.ID, req.NodeId)
		return &pb.AccessResponse{
			Granted:          false,
			LamportTimestamp: timestamp,
		}, nil
	}

	return &pb.AccessResponse{Granted: true,
		LamportTimestamp: timestamp}, nil
}

func (n *Node) isHigherPriority(req1, req2 *Request) bool {
	if req1 == nil || req2 == nil {
		return false
	}
	if req1.Timestamp == req2.Timestamp {
		return req1.NodeID < req2.NodeID
	}
	return req1.Timestamp < req2.Timestamp
}

func (n *Node) ReleaseAccess(ctx context.Context, req *pb.ReleaseRequest) (*pb.ReleaseResponse, error) {

	// Update Lamport clock on release message
	timestamp := n.UpdateLamportClock(req.LamportTimestamp)

	log.Printf("Node %s received release from %s with Lamport timestamp %d",
		n.ID, req.NodeId, req.LamportTimestamp)

	// Remove the released request from queue
	for i := 0; i < len(n.RequestQueue); i++ {
		if n.RequestQueue[i].NodeID == req.NodeId {
			n.RequestQueue = append(n.RequestQueue[:i], n.RequestQueue[i+1:]...)
			break
		}
	}

	// Process next request in queue if we have one
	if len(n.RequestQueue) > 0 && !n.WantCS {
		nextReq := n.RequestQueue[0]
		if client, ok := n.Peers[nextReq.NodeID]; ok {
			go func() {
				_, err := client.RequestAccess(context.Background(), &pb.AccessRequest{
					NodeId:           n.ID,
					LamportTimestamp: n.GetLamportClock(),
				})
				if err != nil {
					log.Printf("Error sending deferred response to %s: %v", nextReq.NodeID, err)
				}
			}()
		}
	}

	return &pb.ReleaseResponse{Acknowledged: true,
		LamportTimestamp: timestamp}, nil
}

func (n *Node) RequestCriticalSection() {
	if n.InCS || n.WantCS {
		return
	}

	n.WantCS = true
	timestamp := n.GetLamportClock()
	n.CurrentRequest = &Request{
		NodeID:    n.ID,
		Timestamp: timestamp,
	}
	n.ResponseCount = 0

	log.Printf("Node %s requesting critical section access with Lamport timestamp %d",
		n.ID, timestamp)

	responses := make(chan bool, len(n.Peers))

	// Request access from all peers
	for peerID, client := range n.Peers {
		go func(id string, c pb.MutexServiceClient) {
			resp, err := c.RequestAccess(context.Background(), &pb.AccessRequest{
				NodeId:           n.ID,
				LamportTimestamp: timestamp,
			})

			if err != nil {
				log.Printf("Error requesting access from %s: %v", id, err)
				responses <- false
				return
			}

			n.UpdateLamportClock(resp.LamportTimestamp)
			responses <- resp.Granted
		}(peerID, client)
	}

	// Wait for all responses
	granted := 0
	needed := len(n.Peers)

	for i := 0; i < needed; i++ {
		if <-responses {
			granted++
		}
	}

	if granted == needed {
		n.InCS = true
		n.ExecuteCriticalSection()
	} else {
		n.WantCS = false
	}
}

func (n *Node) ExecuteCriticalSection() {
	log.Printf("Node %s entering critical section with Lamport timestamp %d",
		n.ID, n.LamportClock)

	// Simulate critical section work
	time.Sleep(2 * time.Second)

	n.InCS = false
	n.WantCS = false
	releaseTimestamp := n.GetLamportClock()

	log.Printf("Node %s leaving critical section with Lamport timestamp %d",
		n.ID, releaseTimestamp)

	// Send release to all peers
	for peerID, client := range n.Peers {
		go func(id string, c pb.MutexServiceClient) {
			_, err := c.ReleaseAccess(context.Background(), &pb.ReleaseRequest{
				NodeId:           n.ID,
				LamportTimestamp: releaseTimestamp,
			})
			if err != nil {
				log.Printf("Error sending release to %s: %v", id, err)
			}
		}(peerID, client)
	}

	// Clear current request and process next in queue
	n.CurrentRequest = nil
	if len(n.RequestQueue) > 0 {
		nextReq := n.RequestQueue[0]
		n.RequestQueue = n.RequestQueue[1:]
		if client, ok := n.Peers[nextReq.NodeID]; ok {
			go func() {
				_, err := client.RequestAccess(context.Background(), &pb.AccessRequest{
					NodeId:           n.ID,
					LamportTimestamp: n.GetLamportClock(),
				})
				if err != nil {
					log.Printf("Error processing next request for %s: %v", nextReq.NodeID, err)
				}
			}()
		}
	}
}
