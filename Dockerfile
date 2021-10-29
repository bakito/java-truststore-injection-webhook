# Build the manager binary
FROM golang:1.17 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY . /workspace
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o truststore-injector-webhook main.go

FROM gcr.io/distroless/static:nonroot
ENTRYPOINT ["/truststore-injector-webhook"]
WORKDIR /opt/go/
COPY --from=builder /workspace/truststore-injector-webhook /opt/go/truststore-injector-webhook
USER 1001:1001
