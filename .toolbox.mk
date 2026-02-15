## toolbox - start
## Generated with https://github.com/bakito/toolbox

## Current working directory
TB_LOCALDIR ?= $(shell which cygpath > /dev/null 2>&1 && cygpath -m $$(pwd) || pwd)
## Location to install dependencies to
TB_LOCALBIN ?= $(TB_LOCALDIR)/bin
$(TB_LOCALBIN):
	if [ ! -e $(TB_LOCALBIN) ]; then mkdir -p $(TB_LOCALBIN); fi

## Tool Binaries
TB_GOLANGCI_LINT ?= $(TB_LOCALBIN)/golangci-lint
TB_GORELEASER ?= $(TB_LOCALBIN)/goreleaser
TB_HELM_DOCS ?= $(TB_LOCALBIN)/helm-docs
TB_SEMVER ?= $(TB_LOCALBIN)/semver

## Tool Versions
TB_GOLANGCI_LINT_VERSION ?= v2.9.0
TB_GORELEASER_VERSION ?= v2.13.3
TB_HELM_DOCS_VERSION ?= v1.14.2
TB_SEMVER_VERSION ?= v1.1.10

## Tool Installer
.PHONY: tb.golangci-lint
tb.golangci-lint: ## Download golangci-lint locally if necessary.
	@test -s $(TB_GOLANGCI_LINT) || \
		GOBIN=$(TB_LOCALBIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(TB_GOLANGCI_LINT_VERSION)
.PHONY: tb.goreleaser
tb.goreleaser: ## Download goreleaser locally if necessary.
	@test -s $(TB_GORELEASER) || \
		GOBIN=$(TB_LOCALBIN) go install github.com/goreleaser/goreleaser/v2@$(TB_GORELEASER_VERSION)
.PHONY: tb.helm-docs
tb.helm-docs: ## Download helm-docs locally if necessary.
	@test -s $(TB_HELM_DOCS) || \
		GOBIN=$(TB_LOCALBIN) go install github.com/norwoodj/helm-docs/cmd/helm-docs@$(TB_HELM_DOCS_VERSION)
.PHONY: tb.semver
tb.semver: ## Download semver locally if necessary.
	@test -s $(TB_SEMVER) || \
		GOBIN=$(TB_LOCALBIN) go install github.com/bakito/semver@$(TB_SEMVER_VERSION)

## Reset Tools
.PHONY: tb.reset
tb.reset:
	@rm -f \
		$(TB_GOLANGCI_LINT) \
		$(TB_GORELEASER) \
		$(TB_HELM_DOCS) \
		$(TB_SEMVER)

## Update Tools
.PHONY: tb.update
tb.update: tb.reset
	toolbox makefile -f $(TB_LOCALDIR)/Makefile \
		github.com/golangci/golangci-lint/v2/cmd/golangci-lint \
		github.com/goreleaser/goreleaser/v2 \
		github.com/norwoodj/helm-docs/cmd/helm-docs \
		github.com/bakito/semver
## toolbox - end
