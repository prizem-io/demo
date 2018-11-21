# Service Mesh Demo

This repo demostrates running Prizem [proxy](https://github.com/prizem-io/proxy)/[control plane](https://github.com/prizem-io/control-plane), [Jaeger](https://www.jaegertracing.io), and [Istio Mixer](https://istio.io) to connect 3 different services together.

### Setup / Initialization

First, verify that you are not running Postgres or any other service on port `5432` or `9000`.  This can be checked using `lsof -i tcp:5432` and `lsof -i tcp:9000`.  Then run the command below.

```bash
./scripts/setup.sh
```

Next, add this entry to your `/etc/hosts`.  The IP can be anything because it is ignored when transparent proxying is initialized.  You can also use any hostname you'd like - just search and replace `api.prizem.io` in `docker-compose.yaml`.

```
172.99.0.2      api.prizem.io
```

### Running the demo

```bash
docker-compose up --build
```

After the first run, the `--build` flag is not necessary.

### Test it!

```bash
curl -s -k -u demo:demo -X GET \
  https://localhost:8000/hello \
  -H 'Accept: application/json; v=1' \
  -H 'Content-Type: application/json; v=1' | jq
```

You can then navigate to [http://localhost:16686](http://localhost:16686 "Jaeger UI") to access the Jaeger UI.