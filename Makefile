# Include toolbox tasks
include ./.toolbox.mk

# Run go golanci-lint
lint: tb.golangci-lint
	$(TB_GOLANGCI_LINT) run --fix

# Run go mod tidy
tidy:
	go mod tidy

# Run tests
test: tidy lint
	go test ./...  -coverprofile=coverage.out
	go tool cover -func=coverage.out

release: tb.semver tb.goreleaser
	@version=$$($(TB_SEMVER)); \
	git tag -s $$version -m"Release $$version"; \
	git push origin $$version
	$(TB_GORELEASER) --clean --parallelism 2

test-release: tb.goreleaser
	$(TB_GORELEASER) --skip=publish --snapshot --clean

helm-docs: tb.helm-docs
	$(TB_HELM_DOCS)

helm-lint: helm-docs
	helm lint ./chart

helm-template:
	helm template ./chart -n jti

