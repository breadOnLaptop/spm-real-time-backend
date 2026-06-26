#!/bin/bash

echo "========================================="
echo " Stopping SPM Cloud Agent"
echo "========================================="

# Kills ALL containers running the spm_agent image
docker rm -f $(docker ps -a -q --filter ancestor=spm_agent) 2>/dev/null || true

echo "Agent successfully stopped."
echo "========================================="
