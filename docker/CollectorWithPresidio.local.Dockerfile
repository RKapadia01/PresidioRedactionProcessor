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

FROM condaforge/mambaforge
WORKDIR /app

COPY ./presidio_grpc_wrapper/requirements.txt /app

RUN apt-get update && apt-get install -y netcat-traditional && rm -rf /var/lib/apt/lists/*

RUN conda install -c conda-forge spacy && \
  python -m spacy download en_core_web_lg && \
  pip install "presidio_analyzer[transformers]" && \
  pip install -r requirements.txt

COPY --from=builder /app/_build/otelcol-presidio ./otel-collector
COPY ./docker/config.yaml .

COPY ./presidio_grpc_wrapper/*.py /app
COPY ./presidio_grpc_wrapper/*.pyi /app

EXPOSE 4317 4318 50051

COPY ./docker/CollectorWithPresidio.entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]