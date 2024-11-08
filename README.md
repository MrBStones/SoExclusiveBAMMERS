# Distributed Mutual Exclusion System

This system implements distributed mutual exclusion using a variant of the Ricart-Agrawala algorithm.

## Requirements

- Go 1.16 or later
- Protocol Buffers compiler
- gRPC

## Building

1. Build the application:
```zsh
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


The implementation uses a version of the Ricart-Agrawala algorithm with Lamport timestamps for distributed mutual exclusion:

1. Each node maintains a Lamport logical clock that:
   - Increments on each local event
   - Updates to max(local_clock, message_timestamp) + 1 on message receipt

2. When a node wants to enter the critical section:
   - It increments its Lamport clock
   - Sends REQUEST messages with the current timestamp to all other nodes
   - Waits for RELEASE ACCESS responses from all other nodes

3. When a node receives a REQUEST:
   - Updates its Lamport clock based on the message timestamp
   - If it's not in the critical section and doesn't want to enter, it sends a RELEASE ACCESS response
   - If it's in the critical section or wants to enter:
     - If its request has a lower timestamp (or equal timestamp but lower node ID), it defers the response
     - Otherwise, it sends a RELEASE ACCESS response

4. When a node leaves the critical section:
   - It sends RELEASE ACCESS messages to all nodes that have deferred requests
   - Each message includes the current Lamport timestamp

The system guarantees both safety and liveness:
- Safety: The Lamport timestamps create a total ordering of requests
- Liveness: Every request eventually gets access due to the FIFO processing of requests

### Example log sequence:
Node1:
```
2024/11/08 11:11:44 Node node1 connected to peer node2 at localhost:5002
2024/11/08 11:11:46 Node node1 connected to peer node3 at localhost:5003
2024/11/08 11:11:47 Node node1 received request from node2 with Lamport timestamp 1
2024/11/08 11:11:47 Node node1 granting node2 access to CS
2024/11/08 11:11:47 Node node1 requesting critical section access with Lamport timestamp 3
2024/11/08 11:11:47 Node node1 is waiting for 2 releases
2024/11/08 11:11:47 Requesting access from node3
2024/11/08 11:11:47 Requesting access from node2
2024/11/08 11:11:47 Finished requesting access from node2
2024/11/08 11:11:47 Node node1 received release from node3 with Lamport timestamp 4
2024/11/08 11:11:47 Finished requesting access from node3
2024/11/08 11:11:48 Node node1 received request from node3 with Lamport timestamp 5
2024/11/08 11:11:48 Node node1 deferring response to node3
2024/11/08 11:11:49 Node node1 received release from node2 with Lamport timestamp 9
2024/11/08 11:11:49 Node node1 entering critical section with Lamport timestamp 12
2024/11/08 11:11:50 Node node1 received request from node2 with Lamport timestamp 10
2024/11/08 11:11:50 Node node1 deferring response to node2
2024/11/08 11:11:51 Node node1 leaving critical section with Lamport timestamp 14
2024/11/08 11:11:51 Node node1 sending 2 Defered responses with Lamport timestamp 14
2024/11/08 11:11:51 Node node1 granting node3 access to CS
2024/11/08 11:11:51 Node node1 granting node2 access to CS
2024/11/08 11:11:52 Node node1 requesting critical section access with Lamport timestamp 15
2024/11/08 11:11:52 Node node1 is waiting for 2 releases
2024/11/08 11:11:52 Requesting access from node3
2024/11/08 11:11:52 Requesting access from node2
2024/11/08 11:11:52 Finished requesting access from node2
2024/11/08 11:11:52 Finished requesting access from node3
2024/11/08 11:11:53 Node node1 received release from node3 with Lamport timestamp 18
2024/11/08 11:11:54 Node node1 received request from node3 with Lamport timestamp 19
2024/11/08 11:11:54 Node node1 deferring response to node3
^Csignal: interrupt
```

Node2:
```
2024/11/08 11:11:44 Node node2 connected to peer node1 at localhost:5001
2024/11/08 11:11:46 Node node2 connected to peer node3 at localhost:5003
2024/11/08 11:11:47 Node node2 requesting critical section access with Lamport timestamp 1
2024/11/08 11:11:47 Node node2 is waiting for 2 releases
2024/11/08 11:11:47 Requesting access from node1
2024/11/08 11:11:47 Requesting access from node3
2024/11/08 11:11:47 Node node2 received release from node1 with Lamport timestamp 2
2024/11/08 11:11:47 Finished requesting access from node1
2024/11/08 11:11:47 Node node2 received release from node3 with Lamport timestamp 2
2024/11/08 11:11:47 Finished requesting access from node3
2024/11/08 11:11:47 Node node2 entering critical section with Lamport timestamp 6
2024/11/08 11:11:47 Node node2 received request from node1 with Lamport timestamp 3
2024/11/08 11:11:47 Node node2 deferring response to node1
2024/11/08 11:11:48 Node node2 received request from node3 with Lamport timestamp 5
2024/11/08 11:11:48 Node node2 deferring response to node3
2024/11/08 11:11:49 Node node2 leaving critical section with Lamport timestamp 9
2024/11/08 11:11:49 Node node2 sending 2 Defered responses with Lamport timestamp 9
2024/11/08 11:11:49 Node node2 granting node1 access to CS
2024/11/08 11:11:49 Node node2 granting node3 access to CS
2024/11/08 11:11:50 Node node2 requesting critical section access with Lamport timestamp 10
2024/11/08 11:11:50 Node node2 is waiting for 2 releases
2024/11/08 11:11:50 Requesting access from node3
2024/11/08 11:11:50 Requesting access from node1
2024/11/08 11:11:50 Finished requesting access from node1
2024/11/08 11:11:50 Finished requesting access from node3
2024/11/08 11:11:51 Node node2 received release from node1 with Lamport timestamp 14
2024/11/08 11:11:52 Node node2 received request from node1 with Lamport timestamp 15
2024/11/08 11:11:52 Node node2 deferring response to node1
2024/11/08 11:11:53 Node node2 received release from node3 with Lamport timestamp 18
2024/11/08 11:11:53 Node node2 entering critical section with Lamport timestamp 19
2024/11/08 11:11:54 Node node2 received request from node3 with Lamport timestamp 19
2024/11/08 11:11:54 Node node2 deferring response to node3
2024/11/08 11:11:55 Node node2 leaving critical section with Lamport timestamp 21
2024/11/08 11:11:55 Node node2 sending 2 Defered responses with Lamport timestamp 21
^Csignal: interrupt
```

Node3:
```
2024/11/08 11:11:45 Node node3 connected to peer node1 at localhost:5001
2024/11/08 11:11:47 Node node3 connected to peer node2 at localhost:5002
2024/11/08 11:11:47 Node node3 received request from node2 with Lamport timestamp 1
2024/11/08 11:11:47 Node node3 granting node2 access to CS
2024/11/08 11:11:47 Node node3 received request from node1 with Lamport timestamp 3
2024/11/08 11:11:47 Node node3 granting node1 access to CS
2024/11/08 11:11:48 Node node3 requesting critical section access with Lamport timestamp 5
2024/11/08 11:11:48 Node node3 is waiting for 2 releases
2024/11/08 11:11:48 Requesting access from node2
2024/11/08 11:11:48 Requesting access from node1
2024/11/08 11:11:48 Finished requesting access from node1
2024/11/08 11:11:48 Finished requesting access from node2
2024/11/08 11:11:49 Node node3 received release from node2 with Lamport timestamp 9
2024/11/08 11:11:50 Node node3 received request from node2 with Lamport timestamp 10
2024/11/08 11:11:50 Node node3 deferring response to node2
2024/11/08 11:11:51 Node node3 received release from node1 with Lamport timestamp 14
2024/11/08 11:11:51 Node node3 entering critical section with Lamport timestamp 16
2024/11/08 11:11:52 Node node3 received request from node1 with Lamport timestamp 15
2024/11/08 11:11:52 Node node3 deferring response to node1
2024/11/08 11:11:53 Node node3 leaving critical section with Lamport timestamp 18
2024/11/08 11:11:53 Node node3 sending 2 Defered responses with Lamport timestamp 18
2024/11/08 11:11:53 Node node3 granting node2 access to CS
2024/11/08 11:11:53 Node node3 granting node1 access to CS
2024/11/08 11:11:54 Node node3 requesting critical section access with Lamport timestamp 19
2024/11/08 11:11:54 Node node3 is waiting for 2 releases
2024/11/08 11:11:54 Requesting access from node2
2024/11/08 11:11:54 Requesting access from node1
2024/11/08 11:11:54 Finished requesting access from node1
2024/11/08 11:11:54 Finished requesting access from node2
^Csignal: interrupt
```

This demonstrates how Lamport timestamps maintain causality and ensure fair ordering of critical section access requests.