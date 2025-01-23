# PresidioRedactionProcessor

[![Docker Image](https://github.com/RKapadia01/PresidioRedactionProcessor/actions/workflows/docker-build-CollectorOnly.yaml/badge.svg?branch=main)](https://github.com/RKapadia01/PresidioRedactionProcessor/actions/workflows/docker-build-CollectorOnly.yaml)
[![Docker Image With Presidio](https://github.com/RKapadia01/PresidioRedactionProcessor/actions/workflows/docker-build-CollectorWithPresidio.yaml/badge.svg?branch=main)](https://github.com/RKapadia01/PresidioRedactionProcessor/actions/workflows/docker-build-CollectorWithPresidio.yaml)

## Quick Start

If you are using the pre-built containers, there are two versions available:

- `rohankapadia/presidioredactioncollector:latest` - This is the latest version of the Processor, without Presidio.
- `rohankapadia/presidioredactioncollector:withpresidio` - This is the latest version of the Processor, with Presidio.

### Using a Presidio-ready Image

To quickly test out the processor, you can run the container that has Presidio baked-in. In this
mode, the default configuration will point to the Presidio service running in the container.
The communication between the Processor and Presidio is done via gRPC.

```bash
docker run --rm -d \
    -p 4318:4318 -p 4317:4317 \
    rohankapadia/presidioredactioncollector:withpresidio
```

And then you can send the telemetry to port 4317/4318. To look at the telemetry ingested, you
can directly look into the logs of the container.

### Deploy your own Presidio Service

It is also possible to deploy your own Presidio service and connect it to the Processor.
In this mode, you need to provide your own configuration file, and the communication between
the Processor and Presidio is done via HTTP.

To do this in Docker, you can run the following commands:

```bash
docker run --rm -d -p 5002:3000 mcr.microsoft.com/presidio-analyzer:latest
docker run --rm -d -p 5001:3000 mcr.microsoft.com/presidio-anonymizer:latest
```

Then you need to edit the configuration file to point to the correct Presidio service. You can do this
by editing the `docker/config.yaml` file:

```yaml
    analyzer_endpoint: http://host.docker.internal:5002/analyze
    anonymizer_endpoint: http://host.docker.internal:5001/anonymize
```

Finally, you can run the Processor container:

```bash
docker run --rm -d \
    -p 4318:4318 -p 4317:4317 \
    -v $(pwd)/docker/config.yaml:/app/config.yaml \
    rohankapadia/presidioredactioncollector:latest
```


## Build the Docker containers:

If for some reason you would like to build the Docker containers with the Collector yourself,
you can do so by running the following commands:

- Build based on local codebase:
    - Collector Only: `docker build -f ./docker/CollectorOnly.local.Dockerfile .`
    - Collector with Presidio: `docker build -f ./docker/CollectorWithPresidio.local.Dockerfile .`
- Build based on published codebase:
    - Collector Only: `docker build -f ./docker/CollectorOnly.Dockerfile .`
    - Collector with Presidio: `docker build -f ./docker/CollectorWithPresidio.Dockerfile .`


## Compile the proto file

If you are making changes to the interface between the Processor and Presidio Wrapper, you would need
to re-compile the gRPC definition to generate the latest interface definition.

You would need the following dependencies:

- Prerequisites for gRPC - Golang: https://grpc.io/docs/languages/go/quickstart/
- Prerequisites for gRPC - Python: https://grpc.io/docs/languages/python/quickstart/

Then you can run the following script to generate the Protocol buffer files:

```bash
# Generate Golang Files
protoc \
    --go_out=./presidioredactionprocessor --go_opt=paths=source_relative \
    --go-grpc_out=./presidioredactionprocessor --go-grpc_opt=paths=source_relative \
    ./presidio.proto

# Generate Python Files
python -m grpc_tools.protoc \
    --proto_path=. \
    --python_out=./presidio_grpc_wrapper \
    --pyi_out=./presidio_grpc_wrapper \
    --grpc_python_out=./presidio_grpc_wrapper \
    ./presidio.proto
```

## Benchmarks

```
loadtest -c 1 -n 1000 -t 10 -P '{"resourceSpans":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my-example-service"}}]},"scopeSpans":[{"scope":{"name":"example-instrumentation-scope","version":"1.0.0"},"spans":[{"traceId":"0123456789abcdef0123456789abcdef","spanId":"0123456789abcdef","name":"my-span-name","kind":"SPAN_KIND_INTERNAL","startTimeUnixNano":"1674457900000000000","endTimeUnixNano":"1674458000000000000","status":{"code":"STATUS_CODE_OK"},"attributes":[{"key":"user.name","value":{"stringValue":"Jacob Zhou"}},{"key":"user.email","value":{"stringValue":"jacob.zhou@example.com"}}]}]}]}]}' -T 'application/json' -m POST http://localhost:4318/v1/traces
```

## Collector -> HTTP * 2 -> Presidio

```
Target URL:          http://localhost:4318/v1/traces
Max requests:        1000
Concurrent clients:  32
Running on cores:    16
Agent:               none

Completed requests:  1000
Total errors:        0
Total time:          7.905 s
Mean latency:        248.8 ms
Effective rps:       127

Percentage of requests served within a certain time
  50%      240 ms
  90%      321 ms
  95%      337 ms
  99%      352 ms
 100%      354 ms (longest request)
```

## Collector -> gRPC * 2 -> Presidio

```
Target URL:          http://localhost:4318/v1/traces
Max requests:        1000
Concurrent clients:  32
Running on cores:    16
Agent:               none

Completed requests:  1000
Total errors:        0
Total time:          8.56 s
Mean latency:        268.8 ms
Effective rps:       117

Percentage of requests served within a certain time
  50%      269 ms
  90%      301 ms
  95%      308 ms
  99%      326 ms
 100%      347 ms (longest request)
```
