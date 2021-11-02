# Build the manager binary
FROM golang:1.17 as builder
RUN apt-get update && apt-get install -y upx

WORKDIR /workspace
# Copy the Go Modules manifests
COPY . /workspace
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o java-truststore-injection-webhook main.go && \
    upx -q java-truststore-injection-webhook

FROM gcr.io/distroless/static:nonroot
ENTRYPOINT ["/opt/go/java-truststore-injection-webhook"]
WORKDIR /opt/go/
COPY --from=builder /workspace/java-truststore-injection-webhook /opt/go/java-truststore-injection-webhook
USER 1001:1001
