FROM golang AS builder

WORKDIR /app

COPY ./docker/builder-config.local.yaml builder-config.yaml
COPY ./presidioredactionprocessor ./presidioredactionprocessor

RUN curl --proto '=https' --tlsv1.2 -fL -o ocb \
  https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv0.117.0/ocb_0.117.0_linux_amd64 && \
  chmod +x ocb

RUN ./ocb --verbose --config builder-config.yaml

FROM python:3.12-slim
WORKDIR /app

RUN pip install presidio_analyzer && \
    pip install presidio_anonymizer && \
    python -m spacy download en_core_web_lg

COPY --from=builder /app/_build/otelcol-presidio ./otel-collector
COPY ./docker/config.yaml .

COPY ./docker/local_scripts/* /

EXPOSE 4317 4318

CMD ["./otel-collector", "--config", "/app/config.yaml"]
