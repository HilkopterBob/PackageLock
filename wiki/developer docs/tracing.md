# Tracing

Tracing is an optional feature for performance analysis.
It's based on the [OpenTelemetry Golang SDK](https://github.com/open-telemetry/opentelemetry-go) and [Jaeger Telemetry](https://www.jaegertracing.io/).

## Prerequisites

1. Docker && Docker Runtime
2. PackageLock Binary

## Usage

1. set `TRACING_ENABLED` to `true` (`$ export TRACING_ENABLED=true`)
2. Start the Jaeger Docker Container:
```bash
docker run -d --name jaeger \ 
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:latest
```
  -  If Jaeger is running somewhere else, specify the location  
      with: `$ export TRACING_JAEGER_URL=http://<ip/domain name>:<port>/api/traces`  

3. start PackageLock (eg. `go run . start`)

PackageLock now traces itself and exports it to the `jaeger` Container.

You can now open the Jaeger web interface at: `http://localhost:16686`
