#!/bin/bash

# Cloud Backend URL
ENDPOINT="spm-real-time-backend.onrender.com:443"

echo "========================================="
echo " SPM Cloud Agent Launcher"
echo "========================================="

echo "[1/4] Stopping any existing agent containers..."
docker rm -f $(docker ps -a -q --filter ancestor=spm_agent) 2>/dev/null || true

echo "[2/4] Building the C++ Agent Docker Image..."
docker build -t spm_agent ./agent

echo "[3/4] Launching the Agent in the background..."
docker run -d --name spm_agent_live spm_agent ./spm_agent $ENDPOINT

echo "[4/4] Success! The agent is running and streaming to the cloud."
echo ""
echo "To view the live agent logs, run:"
echo "  docker logs -f spm_agent_live"
echo "========================================="
