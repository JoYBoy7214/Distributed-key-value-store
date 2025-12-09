# Distributed Key-Value Store in Go

A high-performance, fault-tolerant, and sharded Key-Value store built from scratch in Go. This system implements core distributed systems concepts including **Consistent Hashing**, **Data Replication**, and **Read Failover** without relying on external frameworks.

---

## 🚀 Features

- **Sharding & Partitioning:** Distributes data across multiple nodes using a **Consistent Hashing Ring** with Virtual Nodes to ensure even load distribution.
- **Fault Tolerance:** Implements a **Replication Factor of 3**. Data is written to a primary node and two replicas.
- **High Availability:** Features a **Failover Read Strategy**. If the primary node is offline or returns a 404, the proxy automatically queries the replicas.
- **Thread Safety:** Custom thread-safe storage engine using `sync.RWMutex` to handle concurrent read/write operations.
- **Dynamic Scaling:** Supports adding new nodes to the cluster at runtime via an administrative API.
- **Custom Proxy:** A smart reverse proxy that handles routing, load balancing, and replication logic.

---

## 🛠️ Architecture

The system consists of two main components:

1. **Storage Nodes:** Independent HTTP servers that store key-value pairs in memory. They are unaware of each other.  
2. **Smart Proxy:** The entry point for clients. It manages the Consistent Hash Ring and routes requests to the correct nodes based on the key's hash.

### Algorithm Highlights

- **Hashing:** FNV-1a (32-bit) is used for hashing keys.
- **Consistent Hashing:** Uses a sorted slice of hash values + Binary Search (`O(log N)`) to locate nodes on the ring.
- **Virtual Nodes:** Each physical server is mapped to multiple points on the ring (configurable weight) to prevent data hotspots.

---

## 📦 Installation & Setup

### Prerequisites
- Go 1.18+ installed

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/JoYBoy7214/Distributed-key-value-store.git
cd Distributed-key-value-store
```

### 2️⃣ Start the Storage Nodes

Open **3 separate terminals** and run:

```bash
# Terminal 1
go run ./cmd/node/node.go -port=8081

# Terminal 2
go run ./cmd/node/node.go -port=8082

# Terminal 3
go run ./cmd/node/node.go -port=8083
```

### 3️⃣ Start the Proxy Server

Open a 4th terminal:

```bash
go run ./cmd/proxy/proxy.go
```

- Proxy listens on: **port 8080**

---

## 🔌 API Usage

### 1️⃣ Initialize the Cluster (Add Nodes)

```bash
curl -X POST http://localhost:8080/AddServer \
     -H "Content-Type: application/json" \
     -d '{"url": "http://localhost:8081", "weight": 1}'

curl -X POST http://localhost:8080/AddServer \
     -H "Content-Type: application/json" \
     -d '{"url": "http://localhost:8082", "weight": 1}'

curl -X POST http://localhost:8080/AddServer \
     -H "Content-Type: application/json" \
     -d '{"url": "http://localhost:8083", "weight": 1}'
```

### 2️⃣ Store Data (PUT)

```bash
curl -X PUT http://localhost:8080/PUT \
     -H "Content-Type: application/json" \
     -d '{"Key": "language", "Value": "Go"}'
```

### 3️⃣ Retrieve Data (GET)

```bash
curl -X GET http://localhost:8080/GET \
     -H "Content-Type: application/json" \
     -d '{"Key": "language"}'
```

**Response:**

```json
"Go"
```

---

## 🧪 Testing Fault Tolerance

1️⃣ Start all nodes and proxy  
2️⃣ Add all servers via `/AddServer`  
3️⃣ PUT a value (e.g., key="test")  
4️⃣ Kill the node that owns that key  
5️⃣ GET the key again  

➡ Result: The proxy detects failure and retrieves value from a replica.

---

## 📂 Project Structure

```
├── cmd/
│   ├── node/          # Storage Node entry point
│   └── proxy/         # Proxy/Load Balancer entry point
├── msg/               # Shared message structs (Putmsg, Getmsg)
├── go.mod             # Go module definition
└── README.md          # Project documentation
```

---

## 🔮 Future Improvements

- Data persistence (survive restarts)
- Hinted Handoff to restore missing replica data
- gRPC based communication for lower latency

---
