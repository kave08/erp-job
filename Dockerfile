# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.22.12

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bookworm AS build

WORKDIR /src

ENV CGO_ENABLED=0 \
    GOFLAGS=-trimpath

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY main.go ./
COPY cmd ./cmd
COPY internal ./internal

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -buildvcs=false -ldflags="-s -w" -o /out/erp-job .

FROM gcr.io/distroless/base-debian12:nonroot AS runtime

WORKDIR /app

ENV TZ=UTC

COPY --from=build --chown=nonroot:nonroot /out/erp-job /usr/local/bin/erp-job

ENTRYPOINT ["/usr/local/bin/erp-job"]
CMD ["transfer", "--config-path", "/config/env.yml"]
