#!/bin/bash
curl -X PUT \
  http://control-plane:8000/v1/routes \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "frontend",
  "description": "Demo frontend service",
  "version": {
    "versionLocations": [
      "content-type"
    ],
    "defaultVersion": "v1"
  },
  "authentication": "none",
  "timeout": "10s",
  "operations": [
    {
      "name": "hello",
      "method": "GET",
      "uriPattern": "/hello",
      "timeout": "10s",
      "retry": {
        "attempts": 3,
        "perTryTimeout": "2s"
      },
      "policies": []
    }
  ]
}'

curl -X PUT \
  http://control-plane:8000/v1/routes \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "backend",
  "description": "Demo backend service",
  "version": {
    "versionLocations": [
      "content-type"
    ],
    "defaultVersion": "v1"
  },
  "authentication": "none",
  "timeout": "10s",
  "operations": [
    {
      "name": "message",
      "method": "GET",
      "uriPattern": "/message",
      "timeout": "10s",
      "retry": {
        "attempts": 3,
        "perTryTimeout": "2s"
      },
      "policies": []
    }
  ]
}'

curl -X PUT \
  http://control-plane:8000/v1/routes \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "message",
  "description": "Demo message service",
  "version": {
    "versionLocations": [
      "content-type"
    ],
    "defaultVersion": "v1"
  },
  "authentication": "none",
  "timeout": "10s",
  "operations": [
    {
      "name": "sayHello",
      "method": "GET",
      "uriPattern": "/sayHello",
      "timeout": "10s",
      "retry": {
        "attempts": 3,
        "perTryTimeout": "2s"
      },
      "policies": []
    }
  ]
}'