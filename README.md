# Distributed Mutual Exclusion System

This system implements distributed mutual exclusion using a variant of the Ricart-Agrawala algorithm.

## Requirements

- Go 1.16 or later
- Protocol Buffers compiler
- gRPC

## Building

1. Build the application:
```bash
go build -o mutex
```

## Running

Start multiple nodes with different IDs and addresses:

```zsh
# Start first node
./mutex -id node1 -addr localhost:5001 -peers node2@localhost:5002,node3@localhost:5003

# Start second node
./mutex -id node2 -addr localhost:5002 -peers node1@localhost:5001,node3@localhost:5003

# Start third node
./mutex -id node3 -addr localhost:5003 -peers node1@localhost:5001,node2@localhost:5002
```

Make sure that all peers have correct refferences to eachother or there will be side effects.

## Algorithm Description


The implementation uses a modified version of the Ricart-Agrawala algorithm with Lamport timestamps for distributed mutual exclusion:

1. Each node maintains a Lamport logical clock that:
   - Increments on each local event
   - Updates to max(local_clock, message_timestamp) + 1 on message receipt

2. When a node wants to enter the critical section:
   - It increments its Lamport clock
   - Sends REQUEST messages with the current timestamp to all other nodes
   - Waits for GRANT responses from all nodes

3. When a node receives a REQUEST:
   - Updates its Lamport clock based on the message timestamp
   - If it's not in the critical section and doesn't want to enter, it sends a GRANT
   - If it's in the critical section or wants to enter:
     - If its request has a lower timestamp (or equal timestamp but lower node ID), it defers the response
     - Otherwise, it sends a GRANT

4. When a node leaves the critical section:
   - It sends RELEASE messages to all nodes that have deferred requests
   - Each message includes the current Lamport timestamp

The system guarantees both safety and liveness:
- Safety: The Lamport timestamps create a total ordering of requests
- Liveness: Every request eventually gets access due to the FIFO processing of requests

Example log sequence:
```
Node1: Requesting critical section access with Lamport timestamp 1
Node2: Received request from Node1 with Lamport timestamp 1
Node2: Updated local clock to 2
Node2: Granted access to Node1
Node3: Received request from Node1 with Lamport timestamp 1
Node3: Updated local clock to 2
Node3: Granted access to Node1
Node1: Entering critical section with Lamport timestamp 4
Node1: Leaving critical section with Lamport timestamp 5
```

This demonstrates how Lamport timestamps maintain causality and ensure fair ordering of critical section access requests.


:)