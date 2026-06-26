# SPM: Scalable Performance Monitor

An enterprise-grade, high-performance distributed systems observability stack. SPM collects, aggregates, and visualizes system telemetry across an entire fleet of server nodes in real-time.

## 🚀 Architecture

The system is built on a modern, highly-scalable microservice architecture entirely driven by **Protocol Buffers** and **gRPC**:

1. **C++ Agent (Edge Nodes):** A highly optimized C++ telemetry generator. It collects system-level metrics (CPU, Memory, Disk I/O, Network, Temperature, Top Processes) and streams binary `.proto` payloads over gRPC directly to the central API.
2. **Go Backend (Orchestrator):** A concurrent `golang:latest` backend that listens on a native `tcp:50051` gRPC server. It ingests metrics, writes persistence to PostgreSQL, caches ultra-fast latest states to Redis, and multiplexes the live streams via WebSockets to the frontend.
3. **React Frontend (Dashboard):** An AWS-styled, high-density React web application. It consumes WebSocket streams to render real-time telemetry, segregated process lists, and dynamic filters without requiring page reloads.

## 🛠 Tech Stack
- **Languages:** C++ (Agent), Go (Backend), Javascript/React (Frontend)
- **Transport:** gRPC, Protocol Buffers, WebSockets, REST
- **Storage:** PostgreSQL 15, Redis 7
- **Deployment:** Docker, Docker Compose

## 📦 Running the Stack

The entire stack is containerized for instant deployment. 

```bash
docker-compose up --build -d
```

This spins up 5 containers:
- `spm_db` (Postgres)
- `spm_cache` (Redis)
- `spm_api` (Go API & gRPC Server)
- `spm_frontend` (React + Nginx)
- `spm_agent` (C++ gRPC client)

Once running, access the dashboard at: `http://localhost:3000`

## 📊 Features
- **Real-Time Dashboards:** Millisecond-latency telemetry pushed to the browser.
- **Segregated Process Explorer:** Deep inspection of node-level top processes with dual-filtering.
- **gRPC Health Checks:** Native gRPC Health checking validating PostgreSQL and Redis connections.
- **Data Export:** Instant CSV export of historical telemetry metrics.
- **Dynamic Provisioning:** Easily scale the C++ agents horizontally.
