docker network create --driver bridge mesh-network

# Postgres
docker run --rm --name postgres1 \
    --network mesh-network \
    -v $(pwd)/pgdata:/var/lib/postgresql/data \
    -p 5432:5432 \
    -e POSTGRES_PASSWORD=mysecretpassword \
    -d postgres
./migrate/standup.sh

# Control Plane
docker run --rm --name control-plane1 \
    --network mesh-network \
    -p 9000:8000 \
    -v $(pwd)/etc/control-plane:/app/etc \
    -t prizem-io/control-plane:latest

# Register the services
docker run --rm --name register1 \
    --network mesh-network \
    -t prizem-io/demo-register:latest

# Proxy/sidecar (for all services)
docker run --rm --name proxy1 \
    --network mesh-network \
    -v $(pwd)/etc/proxy:/app/etc \
    -p 8000:50052 \
    -p 6060:6060 \
    -e REST_API_URI=http://control-plane1:8000 \
    -e GRPC_API_TARGET=control-plane1:9000 \
    -e MIXER_TARGET=mixer1:9091 \
    -t prizem-io/proxy:latest

# Frontend
docker run --rm --name frontend1 \
    --network mesh-network \
    -e BACKEND_URI=http://proxy1:50062/message \
    -t prizem-io/demo-frontend:latest

# Backend
docker run --rm --name backend1 \
    --network mesh-network \
    -e MESSAGE_URI=http://proxy1:50062/sayHello \
    -t prizem-io/demo-backend:latest

# Message service
docker run --rm --name message1 \
    --network mesh-network \
    -t prizem-io/demo-message:latest

# TODO - allow for hostnames in addition to IP addresses
docker run --rm --name announce1 \
    --network mesh-network \
    -t prizem-io/demo-announce:latest

# Transparent proxying - Work In Progress
docker run --rm --name proxy-init1 \
    --privileged \
    --net=host \
    --cap-add=NET_ADMIN \
    --cap-add=NET_RAW \
    -e NODE_NAME=proxy1 \
    -t prizem-io/demo-proxy-init:latest -p 50062 -s PRIZEM -m false

#     --network mesh-network \

# Get IP Addresses for each service
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' proxy1
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' frontend1
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' backend1
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' message1

# Manually annouce services for now
curl -H "X-Target-Host: `docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' frontend1`" -d '{"id": "f275d3a9-fe9a-4aee-acfa-7a1e51a95c71", "service": "frontend", "name": "frontend-01", "ports": [{"port": 8000, "protocol": "HTTP/1"}]}' -X POST http://localhost:6060/register
curl -H "X-Target-Host: `docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' backend1`" -d '{"id": "f69670d3-d388-4681-84c4-a1173ab06636", "service": "backend", "name": "backend-01", "ports": [{"port": 8000, "protocol": "HTTP/1"}]}' -X POST http://localhost:6060/register
curl -H "X-Target-Host: `docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' message1`" -d '{"id": "3b52cf50-791b-4c05-82f2-c59d349ab889", "service": "message", "name": "message-01", "ports": [{"port": 8000, "protocol": "HTTP/1"}]}' -X POST http://localhost:6060/register


curl -H "X-Target-Host: `docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' demo_frontend_1`" -d '{"id": "f275d3a9-fe9a-4aee-acfa-7a1e51a95c71", "service": "frontend", "name": "frontend-01", "ports": [{"port": 8000, "protocol": "HTTP/1"}]}' -X POST http://localhost:6060/register
curl -H "X-Target-Host: `docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' demo_backend_1`" -d '{"id": "f69670d3-d388-4681-84c4-a1173ab06636", "service": "backend", "name": "backend-01", "ports": [{"port": 8000, "protocol": "HTTP/1"}]}' -X POST http://localhost:6060/register
curl -H "X-Target-Host: `docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' demo_message_1`" -d '{"id": "3b52cf50-791b-4c05-82f2-c59d349ab889", "service": "message", "name": "message-01", "ports": [{"port": 8000, "protocol": "HTTP/1"}]}' -X POST http://localhost:6060/register

# Test it!
curl -s -k -u demo:demo -X GET \
  https://localhost:8000/hello \
  -H 'Accept: application/json; v=1' \
  -H 'Content-Type: application/json; v=1' | jq


echo "
rdr pass inet proto tcp from any to any port 80 -> `docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' demo_proxy_1` port 50062
" | sudo pfctl -ef -
