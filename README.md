# Service Mesh Demo

This repo demostrates running Prizem [proxy](https://github.com/prizem-io/proxy)/[control plane](https://github.com/prizem-io/control-plane), [Jaeger](https://www.jaegertracing.io), and [Istio Mixer](https://istio.io) to connect 3 different services together.

### Initialization

```bash
./scripts/setup.sh
```

### Running the demo

```bash
docker-compose up
```

### Test it!

```bash
curl -s -k -u demo:demo -X GET \
  https://localhost:8000/hello \
  -H 'Accept: application/json; v=1' \
  -H 'Content-Type: application/json; v=1' | jq
```

You can then navigate to [http://localhost:16686](http://localhost:16686 "Jaeger UI") to access the Jaeger UI.