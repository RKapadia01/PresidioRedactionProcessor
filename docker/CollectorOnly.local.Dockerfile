FROM golang AS builder

WORKDIR /app

COPY ./docker/builder-config.local.yaml builder-config.yaml
COPY ./presidioredactionprocessor ./presidioredactionprocessor

RUN curl --proto '=https' --tlsv1.2 -fL -o ocb \
  https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv0.117.0/ocb_0.117.0_linux_amd64 && \
  chmod +x ocb

ENV GO111MODULE=on
ENV GOMAXPROCS=4
ENV GOOS=linux
ENV GOARCH=amd64

RUN ./ocb --verbose --config builder-config.yaml

FROM debian:stable-slim
WORKDIR /app

COPY --from=builder /app/_build/otelcol-presidio ./otel-collector
COPY ./docker/config.yaml .

EXPOSE 4317 4318

ENTRYPOINT ["./otel-collector"]
CMD ["--config", "/app/config.yaml"]
