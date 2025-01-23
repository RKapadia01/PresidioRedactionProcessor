# PresidioRedactionProcessor

## Quick Start

To quickly test out the processor, you can run the container that has Presidio baked-in:

```bash
docker run --rm -d \
    -p 4318:4318 \
    -p 4317:4317 \
    -p 50051:50051 \
    rohankapadia/presidioredactioncollector:withpresidio
```

And then you can send the telemetry to port 4317/4318. To look at the telemetry ingested, you
can directly look into the logs of the container.


## Build the Docker containers:

- Local version: `docker build -f ./docker/CollectorOnly.local.Dockerfile .`
- Published version: `docker build -f ./docker/CollectorOnly.Dockerfile .`
- Local version with Presidio: `docker build -f ./docker/CollectorWithPresidio.local.Dockerfile .`
- Published version with Presidio: `docker build -f ./docker/CollectorWithPresidio.Dockerfile .`


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
