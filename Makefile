# Include toolbox tasks
include ./.toolbox.mk

# Run go golanci-lint
lint: golangci-lint
	$(GOLANGCI_LINT) run --fix

# Run go mod tidy
tidy:
	go mod tidy

# Run tests
test: tidy lint
	go test ./...  -coverprofile=coverage.out
	go tool cover -func=coverage.out

# Run tests
release: goreleaser semver
	@version=$$($(LOCALBIN)/semver); \
	git tag -s $$version -m"Release $$version"
	$(GORELEASER) --clean

test-release: goreleaser
	$(GORELEASER) --skip=publish --snapshot --clean

docs: helm-docs
	$(HELM_DOCS)

helm-lint: docs
	helm lint ./chart

helm-template:
	helm template ./chart -n jti

