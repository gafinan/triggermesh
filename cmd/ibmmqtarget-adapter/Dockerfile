# syntax=docker/dockerfile:1.3

FROM golang:bullseye as builder

RUN apt-get update && \
    apt-get install -y curl && \
    curl https://public.dhe.ibm.com/ibmdl/export/pub/software/websphere/messaging/mqdev/redist/9.2.3.0-IBM-MQC-Redist-LinuxX64.tar.gz -o mq.tar.gz          && \
    mkdir -p /opt/mqm && \
    tar -C /opt/mqm -xzf mq.tar.gz

ENV MQ_INSTALLATION_PATH="/opt/mqm"
ENV CGO_LDFLAGS_ALLOW="-Wl,-rpath.*"
ENV CGO_LDFLAGS="-L$MQ_INSTALLATION_PATH/lib64 -Wl,-rpath,$MQ_INSTALLATION_PATH/lib64"
ENV CGO_CFLAGS="-I$MQ_INSTALLATION_PATH/inc"

WORKDIR /go/triggermesh

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux GOARCH=amd64 go build -v -o ibmmqtarget-adapter ./cmd/ibmmqtarget-adapter 


FROM gcr.io/distroless/base:nonroot

WORKDIR /opt/mqm/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libstdc++.so.6 /usr/lib/x86_64-linux-gnu/
COPY --from=builder /lib/x86_64-linux-gnu/libgcc_s.so.1 /lib/x86_64-linux-gnu/
COPY --from=builder /opt/mqm .
COPY --from=builder /go/triggermesh/ibmmqtarget-adapter .

ENTRYPOINT ["./ibmmqtarget-adapter"]
