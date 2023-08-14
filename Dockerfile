FROM golang:1.21 as build

ENV CGO_ENABLED 1
RUN apt-get update && apt-get install -y libolm-dev && \
  rm -rf /var/lib/apt/lists/*

COPY . /app

ARG VERSION
RUN cd /app && \
  go build -ldflags="-s -w -X main.version=$VERSION" -trimpath -o /matrix-gpt ./cmd/matrix-gpt

FROM ubuntu:22.04
RUN apt-get update && \
  apt-get install -y libolm3 ca-certificates tzdata && \
  rm -rf /var/lib/apt/lists/*

COPY --from=build /matrix-gpt /matrix-gpt
USER 1337
CMD ["/matrix-gpt"]
