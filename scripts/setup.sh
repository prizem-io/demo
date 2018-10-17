#!/bin/bash

# Run Postgres
echo "Running Postgres..."
docker run --rm --name postgres \
    -v $(pwd)/pgdata:/var/lib/postgresql/data \
    -p 5432:5432 \
    -e POSTGRES_PASSWORD=mysecretpassword \
    -d postgres

echo "Waiting for Postgres to start..."
sleep 10

# Initialize database schema
echo "Initializing database schema..."
./migrate/standup.sh

# Run Control Plane
echo "Running control plane..."
docker run --rm -d --name control-plane \
    -p 9000:8000 \
    -v $(pwd)/etc/control-plane:/app/etc \
    -t prizem-io/control-plane:latest

echo "Waiting for control-plane to start..."
sleep 10

# Register demo services
echo "Registering services..."
./scripts/register.sh

# Cleanup
echo "Shutting down..."
docker stop control-plane postgres
