#!/usr/bin/env bash
set -e

# Start Processor A in the background
echo "Starting Presidio Wrapper..."
python /app/server.py &
PROCESSOR_A_PID=$!

# Start Processor B in the background
echo "Starting OTLP Collector..."
./otel-collector --config /app/config.yaml &
PROCESSOR_B_PID=$!

# Wait for any process to exit
wait -n

# Exit with the status of the process that exited first
exit $?