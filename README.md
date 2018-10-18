# Service Mesh Demo

This repo demostrates running Prizem [proxy](https://github.com/prizem-io/proxy)/[control plane](https://github.com/prizem-io/control-plane), [Jaeger](https://www.jaegertracing.io), and [Istio Mixer](https://istio.io) to connect 3 different services together.

### Setup / Initialization

```bash
./scripts/setup.sh
```

Then add this entry to your `/etc/hosts`.  Really the IP can be anything because it is ignored when transparent proxying is initialized.  You can also use whatever hostname you'd - just search and replace `api.prizem.io` in `docker-compose.yaml`.

```
172.99.0.2      api.prizem.io
```

### Running the demo

```bash
docker-compose up --build
```

### Test it!

```bash
curl -s -k -u demo:demo -X GET \
  https://localhost:8000/hello \
  -H 'Accept: application/json; v=1' \
  -H 'Content-Type: application/json; v=1' | jq
```

You can then navigate to [http://localhost:16686](http://localhost:16686 "Jaeger UI") to access the Jaeger UI.