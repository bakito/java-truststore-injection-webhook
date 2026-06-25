# Build the manager binary
FROM golang:1.26-alpine@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 AS builder
RUN apk update && apk add upx

WORKDIR /workspace
# Copy the Go Modules manifests
COPY . /workspace
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o java-truststore-injection-webhook main.go && \
    upx -q java-truststore-injection-webhook

FROM gcr.io/distroless/static:nonroot@sha256:963fa6c544fe5ce420f1f54fb88b6fb01479f054c8056d0f74cc2c6000df5240
ENTRYPOINT ["/opt/go/java-truststore-injection-webhook"]
WORKDIR /opt/go/
COPY --from=builder /workspace/java-truststore-injection-webhook /opt/go/java-truststore-injection-webhook
USER 1001:1001
