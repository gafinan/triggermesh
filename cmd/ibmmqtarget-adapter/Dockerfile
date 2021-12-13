FROM golang:buster

RUN apt-get update && \
    apt-get install -y curl && \
    curl https://public.dhe.ibm.com/ibmdl/export/pub/software/websphere/messaging/mqdev/redist/9.2.3.0-IBM-MQC-Redist-LinuxX64.tar.gz -o mq.tar.gz          && \
    mkdir -p /opt/mqm             && \
    tar -C /opt/mqm -xzf mq.tar.gz

ENV MQ_INSTALLATION_PATH="/opt/mqm"
ENV CGO_LDFLAGS_ALLOW="-Wl,-rpath.*"
ENV CGO_LDFLAGS="-L$MQ_INSTALLATION_PATH/lib64 -Wl,-rpath,$MQ_INSTALLATION_PATH/lib64"
ENV CGO_CFLAGS="-I$MQ_INSTALLATION_PATH/inc"

WORKDIR /go/ibm-mq-target
COPY . .

RUN go build -v -o ibmmqtarget-adapter ./cmd/ibmmqtarget-adapter 


FROM debian:buster-slim

WORKDIR /opt/mqm/
COPY --from=0 /opt/mqm .
COPY --from=0 /go/ibm-mq-target/ibmmqtarget-adapter .

ENTRYPOINT ["./ibmmqtarget-adapter"]
