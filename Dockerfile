FROM --platform=$BUILDPLATFORM golang:1.24 AS builder

RUN apt-get update && apt-get install -y gcc-aarch64-linux-gnu

WORKDIR /opt/service

ENV GO111MODULE=on

COPY . ./

RUN go generate -C cmd
RUN CGO_ENABLED=1 CGO_CFLAGS="-D_LARGEFILE64_SOURCE" go test -mod vendor ./...

ARG TARGETOS TARGETARCH

RUN if [ "$TARGETARCH" = "arm64" ] ; then export CC=aarch64-linux-gnu-gcc; fi && \
GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=1 CGO_CFLAGS="-D_LARGEFILE64_SOURCE" go build -C cmd -mod vendor \
-ldflags="-X 'main.Version=$(git rev-parse --verify --short HEAD)' -X \"main.BuildAt=$(date '+%d.%m.%Y %H:%M:%S')\" -extldflags=-static" \
-o ../bin/svc

# вторая ступень
FROM alpine:3.20

WORKDIR /opt/service

COPY --from=builder /opt/service/bin/. .
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV LOG_LEVEL=debug

ENTRYPOINT ["/opt/service/svc"]
