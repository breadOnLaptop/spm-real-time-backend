#!/bin/bash

echo "========================================="
echo " Stopping SPM Cloud Agent"
echo "========================================="

docker rm -f spm_agent_live 2>/dev/null || true

echo "Agent successfully stopped."
echo "========================================="
