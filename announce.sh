#!/bin/bash
curl -d '{"id": "f275d3a9-fe9a-4aee-acfa-7a1e51a95c71", "service": "frontend", "name": "frontend-01", "ports": [{"port": 8000, "protocol": "HTTP/1"}]}' -X POST http://proxy1:6060/register
curl -d '{"id": "f69670d3-d388-4681-84c4-a1173ab06636", "service": "backend", "name": "backend-01", "ports": [{"port": 8000, "protocol": "HTTP/1"}]}' -X POST http://proxy1:6060/register
curl -d '{"id": "3b52cf50-791b-4c05-82f2-c59d349ab889", "service": "message", "name": "message-01", "ports": [{"port": 8000, "protocol": "HTTP/1"}]}' -X POST http://proxy1:6060/register