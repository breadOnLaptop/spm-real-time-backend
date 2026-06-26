#!/bin/bash

echo "Stopping existing containers..."
docker-compose down

echo "Starting Docker Compose services (API, DB, Redis, Frontend, Agent)..."
docker-compose up --build -d

echo "=========================================="
echo "SPM MVP is up and running!"
echo "Dashboard: http://localhost:3000"
echo "API Backend: http://localhost:8000"
echo "All components (including the C++ Agent) are now fully containerized in Docker!"
echo "To stop everything, run: docker-compose down"
echo "=========================================="
