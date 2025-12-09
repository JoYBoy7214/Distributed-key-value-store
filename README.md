# Distributed Key-Value Store in Go

A high-performance, fault-tolerant, and sharded Key-Value store built from scratch in Go. This system implements core distributed systems concepts including **Consistent Hashing**, **Data Replication**, and **Read Failover** without relying on external frameworks.

## 🚀 Features

* **Sharding & Partitioning:** Distributes data across multiple nodes using a **Consistent Hashing Ring** with Virtual Nodes to ensure even load distribution.
* **Fault Tolerance:** Implements a **Replication Factor of 3**. Data is written to a primary node and two replicas.
* **High Availability:** Features a **Failover Read Strategy**. If the primary node is offline or returns a 404, the proxy automatically queries the replicas.
* **Thread Safety:** Custom thread-safe storage engine using `sync.RWMutex` to handle concurrent read/write operations.
* **Dynamic Scaling:** Supports adding new nodes to the cluster runtime via an administrative API.
* **Custom Proxy:** A smart reverse proxy that handles routing, load balancing, and replication logic.

## 🛠️ Architecture

The system consists of two main components:

1.  **Storage Nodes:** Independent HTTP servers that store key-value pairs in memory. They are unaware of each other.
2.  **Smart Proxy:** The entry point for clients. It manages the Consistent Hash Ring and routes requests to the correct nodes based on the key's hash.

### Algorithm Highlights
* **Hashing:** FNV-1a (32-bit) is used for hashing keys.
* **Consistent Hashing:** Uses a sorted slice of hash values + Binary Search (`O(log N)`) to locate nodes on the ring.
* **Virtual Nodes:** Each physical server is mapped to multiple points on the ring (configurable weight) to prevent data hotspots.

## 📦 Installation & Setup

### Prerequisites
* Go 1.18+ installed

### 1. Clone the Repository
```bash
git clone [https://github.com/JoYBoy7214/Distributed-key-value-store.git](https://github.com/JoYBoy7214/Distributed-key-value-store.git)
cd distributed-kv-store