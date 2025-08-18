# Build the manager binary
FROM golang:1.25-alpine AS builder
RUN apk update && apk add upx

WORKDIR /workspace
# Copy the Go Modules manifests
COPY . /workspace
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o java-truststore-injection-webhook main.go && \
    upx -q java-truststore-injection-webhook

FROM gcr.io/distroless/static:nonroot
ENTRYPOINT ["/opt/go/java-truststore-injection-webhook"]
WORKDIR /opt/go/
COPY --from=builder /workspace/java-truststore-injection-webhook /opt/go/java-truststore-injection-webhook
USER 1001:1001
