package peer

import (
	"container/list"
	"context"
	"fmt"
	"log"
	pb "mutex/stc"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Node struct {
	ID                string
	Address           string
	Peers             map[string]pb.MutexServiceClient
	LamportClock      uint64
	WantCS            bool // state == WANTED;
	InCS              bool // state == HELD; state == RELEASED if neither is true;
	DeferredResponses *list.List
	ResponseCount     int
	CurrentRequest    *pb.AccessRequest
	ReqMu             sync.Mutex
	LamMu             sync.Mutex
	CsMu              sync.Mutex
	Release           (chan bool)
	pb.UnimplementedMutexServiceServer
}

// --- initialization functions ---
func NewNode(id, address string) *Node {
	return &Node{
		ID:                id,
		Address:           address,
		Peers:             make(map[string]pb.MutexServiceClient),
		LamportClock:      0,
		DeferredResponses: list.New(),
		Release:           make(chan bool, 1),
	}
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

// --- Server functions ---
func (n *Node) RequestAccess(ctx context.Context, req *pb.AccessRequest) (*pb.AccessResponse, error) {
	n.ReqMu.Lock()
	defer n.ReqMu.Unlock()

	// Update Lamport clock on message receipt
	timestamp := n.UpdateLamportClock(req.LamportTimestamp)

	log.Printf("Node %s received request from %s with Lamport timestamp %d", n.ID, req.NodeId, req.LamportTimestamp)

	if n.InCS || (n.WantCS && n.isHigherPriority(n.CurrentRequest, req)) {
		n.DeferredResponses.PushBack(req.NodeId)
		log.Printf("Node %s deferring response to %s", n.ID, req.NodeId)
		return &pb.AccessResponse{Granted: false, LamportTimestamp: timestamp}, nil
	}

	defer n.SendReleaseMSG(req.NodeId, n.Peers[req.NodeId], timestamp) // actual message that matters
	return &pb.AccessResponse{Granted: false, LamportTimestamp: timestamp}, nil
}

func (n *Node) ReleaseAccess(ctx context.Context, req *pb.ReleaseRequest) (*pb.ReleaseResponse, error) {

	// Update Lamport clock on release message
	timestamp := n.UpdateLamportClock(req.LamportTimestamp)

	log.Printf("Node %s received release from %s with Lamport timestamp %d", n.ID, req.NodeId, req.LamportTimestamp)
	n.Release <- true

	return &pb.ReleaseResponse{Acknowledged: true, LamportTimestamp: timestamp}, nil
}

// --- client functions ---
func (n *Node) RequestCriticalSection() {
	n.CsMu.Lock()
	defer n.CsMu.Unlock()
	if n.InCS || n.WantCS {
		return
	}

	n.WantCS = true
	timestamp := n.GetLamportClock()
	n.CurrentRequest = &pb.AccessRequest{
		NodeId:           n.ID,
		LamportTimestamp: timestamp,
	}
	n.ResponseCount = 0

	log.Printf("Node %s requesting critical section access with Lamport timestamp %d", n.ID, timestamp)

	responses := make(chan bool, len(n.Peers))

	// Request access from all peers
	for peerID, client := range n.Peers {
		go func(id string, c pb.MutexServiceClient) {
			log.Printf("Requesting access from %s", peerID)
			resp, err := c.RequestAccess(context.Background(), n.CurrentRequest)

			if err != nil {
				log.Printf("Error requesting access from %s: %v", id, err)
				responses <- false
				return
			}
			log.Printf("Finished requesting access from %s", peerID)
			n.UpdateLamportClock(resp.LamportTimestamp)
			responses <- <-n.Release

		}(peerID, client)
	}

	// Wait for all responses
	log.Printf("Node %s is waiting for %d releases", n.ID, len(n.Peers))
	granted := 0
	needed := len(n.Peers)

	for i := 0; i < needed; i++ {
		if <-responses {
			granted++
		}
	}

	n.ExecuteCriticalSection()

}

func (n *Node) ExecuteCriticalSection() {

	log.Printf("Node %s entering critical section with Lamport timestamp %d", n.ID, n.LamportClock)
	n.InCS = true
	// Simulate critical section work
	time.Sleep(2 * time.Second)

	n.InCS = false
	n.WantCS = false
	releaseTimestamp := n.GetLamportClock()

	log.Printf("Node %s leaving critical section with Lamport timestamp %d", n.ID, releaseTimestamp)

	// Send release to all peers in defered
	log.Printf("Node %s sending %d Defered responses with Lamport timestamp %d", n.ID, n.DeferredResponses.Len(), releaseTimestamp)
	for n.DeferredResponses.Len() > 0 {
		front := n.DeferredResponses.Front()

		peerID := front.Value.(string)
		client := n.Peers[peerID]
		n.SendReleaseMSG(peerID, client, releaseTimestamp)

		n.DeferredResponses.Remove(front)
	}
}

// --- util functions ---
func (n *Node) UpdateLamportClock(msgTimestamp uint64) uint64 {
	n.LamMu.Lock()
	// Same as: clock = max(local_clock, msg_timestamp) + 1
	if msgTimestamp > n.LamportClock {
		n.LamportClock = msgTimestamp
	}
	n.LamportClock++
	n.LamMu.Unlock()
	return n.LamportClock
}

func (n *Node) GetLamportClock() uint64 {
	n.LamMu.Lock()
	n.LamportClock++
	n.LamMu.Unlock()
	return n.LamportClock
}

func (n *Node) isHigherPriority(req1, req2 *pb.AccessRequest) bool {
	if req1 == nil || req2 == nil {
		return false
	}
	if req1.LamportTimestamp == req2.LamportTimestamp {
		return req1.NodeId < req2.NodeId
	}
	return req1.LamportTimestamp < req2.LamportTimestamp
}

func (n *Node) SendReleaseMSG(peerID string, client pb.MutexServiceClient, releaseTimestamp uint64) {
	_, err := client.ReleaseAccess(context.Background(), &pb.ReleaseRequest{
		NodeId:           n.ID,
		LamportTimestamp: releaseTimestamp,
	})
	if err != nil {
		log.Printf("Error sending release to %s: %v", peerID, err)
	}
	log.Printf("Node %s granting %s access to CS", n.ID, peerID)
}
