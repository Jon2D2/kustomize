# Copyright 2019 The Kubernetes Authors.
# SPDX-License-Identifier: Apache-2.0
#
# Makefile for kustomize CLI and API.

MYGOBIN := $(shell go env GOPATH)/bin
TOOLSBIN := $(PWD)/hack/tools/bin
PATH := $(TOOLSBIN):$(MYGOBIN):$(PATH)
SHELL := /bin/bash

.DEFAULT_GOAL := all

export GO111MODULE=on

.PHONY: all
all: verify-kustomize

.PHONY: verify-kustomize
verify-kustomize: \
	lint-kustomize \
	test-unit-kustomize-all \
	test-examples-kustomize-against-HEAD \
	test-examples-kustomize-against-latest

# The following target referenced by a file in
# https://github.com/kubernetes/test-infra/tree/master/config/jobs/kubernetes-sigs/kustomize
.PHONY: prow-presubmit-check
prow-presubmit-check: \
	lint-kustomize \
	test-unit-kustomize-all \
	test-unit-cmd-all \
	test-go-mod

.PHONY: verify-kustomize-e2e
verify-kustomize-e2e: test-examples-e2e-kustomize

# Other builds in this repo might want a different linter version.
# Without one Makefile to rule them all, the different makes
# cannot assume that golanci-lint is at the version they want
# since everything uses the same implicit GOPATH.
# This installs in a temp dir to avoid overwriting someone else's
# linter, then installs in MYGOBIN with a new name.
# Version pinned by hack/tools/go.mod
$(TOOLSBIN)/golangci-lint-kustomize:
	( \
		set -e; \
		cd hack/tools; \
		GO111MODULE=on GOBIN=$(TOOLSBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint; \
		mv $(TOOLSBIN)/golangci-lint $(TOOLSBIN)/golangci-lint-kustomize \
	)

# Version pinned by hack/tools/go.mod
$(TOOLSBIN)/mdrip:
	cd hack/tools; \
	GOBIN=$(TOOLSBIN) go install github.com/monopole/mdrip

# Version pinned by hack/tools/go.mod
$(TOOLSBIN)/stringer:
	cd hack/tools; \
	GOBIN=$(TOOLSBIN) go install golang.org/x/tools/cmd/stringer

# Version pinned by hack/tools/go.mod
$(TOOLSBIN)/goimports:
	cd hack/tools; \
	GOBIN=$(TOOLSBIN) go install golang.org/x/tools/cmd/goimports

# Install resource from whatever is checked out.
$(MYGOBIN)/resource:
	cd cmd/resource; \
	go install .

# To pin pluginator, use this recipe instead:
#	cd hack/tools; \
#	GOBIN=$(TOOLSBIN) go install sigs.k8s.io/kustomize/pluginator/v2
$(TOOLSBIN)/pluginator:
	cd pluginator; \
	GOBIN=$(TOOLSBIN) go install .

# Install kustomize from whatever is checked out.
$(MYGOBIN)/kustomize:
	cd kustomize; \
	GOBIN=$(TOOLSBIN) go install .

.PHONY: install-tools
install-tools: \
	$(TOOLSBIN)/goimports \
	$(TOOLSBIN)/golangci-lint-kustomize \
	$(TOOLSBIN)/mdrip \
	$(TOOLSBIN)/pluginator \
	$(TOOLSBIN)/stringer

### Begin kustomize plugin rules.
#
# The rules to deal with builtin plugins are a bit
# complicated because
#
# - Every builtin plugin is a Go plugin -
#   meaning it gets its own module directory
#   (outside of the api module) with Go
#   code in a 'main' package per Go plugin rules.
# - kustomize locates plugins using the
#   'apiVersion' and 'kind' fields from the
#   plugin config file.
# - k8s wants CamelCase in 'kind' fields.
# - The module name (the last name in the path)
#   must be the lowercased 'kind' of the
#   plugin because Go and related tools
#   demand lowercase in import paths, but
#   allow CamelCase in file names.
# - the generated code must live in the api
#   module (it's linked into the api).

# Where all generated builtin plugin code should go.
pGen=api/builtins
# Where the builtin Go plugin modules live.
pSrc=plugin/builtin

_builtinplugins = \
	AnnotationsTransformer.go \
	ConfigMapGenerator.go \
	HashTransformer.go \
	ImageTagTransformer.go \
	InventoryTransformer.go \
	LabelTransformer.go \
	LegacyOrderTransformer.go \
	NamespaceTransformer.go \
	PatchJson6902Transformer.go \
	PatchStrategicMergeTransformer.go \
	PatchTransformer.go \
	PrefixSuffixTransformer.go \
	ReplicaCountTransformer.go \
	SecretGenerator.go

# Maintaining this explicit list of generated files, and
# adding it as a dependency to a few targets, to assure
# they get recreated if deleted.  The rules below on how
# to make them don't, by themselves, assure they will be
# recreated if deleted.
builtinplugins = $(patsubst %,$(pGen)/%,$(_builtinplugins))

# These rules are verbose, but assure that if a source file
# is modified, the corresponding generated file, and only
# that file, will be recreated.
$(pGen)/AnnotationsTransformer.go: $(pSrc)/annotationstransformer/AnnotationsTransformer.go
$(pGen)/ConfigMapGenerator.go: $(pSrc)/configmapgenerator/ConfigMapGenerator.go
$(pGen)/HashTransformer.go: $(pSrc)/hashtransformer/HashTransformer.go
$(pGen)/ImageTagTransformer.go: $(pSrc)/imagetagtransformer/ImageTagTransformer.go
$(pGen)/InventoryTransformer.go: $(pSrc)/inventorytransformer/InventoryTransformer.go
$(pGen)/LabelTransformer.go: $(pSrc)/labeltransformer/LabelTransformer.go
$(pGen)/LegacyOrderTransformer.go: $(pSrc)/legacyordertransformer/LegacyOrderTransformer.go
$(pGen)/NamespaceTransformer.go: $(pSrc)/namespacetransformer/NamespaceTransformer.go
$(pGen)/PatchJson6902Transformer.go: $(pSrc)/patchjson6902transformer/PatchJson6902Transformer.go
$(pGen)/PatchStrategicMergeTransformer.go: $(pSrc)/patchstrategicmergetransformer/PatchStrategicMergeTransformer.go
$(pGen)/PatchTransformer.go: $(pSrc)/patchtransformer/PatchTransformer.go
$(pGen)/PrefixSuffixTransformer.go: $(pSrc)/prefixsuffixtransformer/PrefixSuffixTransformer.go
$(pGen)/ReplicaCountTransformer.go: $(pSrc)/replicacounttransformer/ReplicaCountTransformer.go
$(pGen)/SecretGenerator.go: $(pSrc)/secretgenerator/SecretGenerator.go

# The (verbose but portable) Makefile way to convert to lowercase.
toLowerCase = $(subst A,a,$(subst B,b,$(subst C,c,$(subst D,d,$(subst E,e,$(subst F,f,$(subst G,g,$(subst H,h,$(subst I,i,$(subst J,j,$(subst K,k,$(subst L,l,$(subst M,m,$(subst N,n,$(subst O,o,$(subst P,p,$(subst Q,q,$(subst R,r,$(subst S,s,$(subst T,t,$(subst U,u,$(subst V,v,$(subst W,w,$(subst X,x,$(subst Y,y,$(subst Z,z,$1))))))))))))))))))))))))))

$(pGen)/%.go: $(TOOLSBIN)/pluginator $(TOOLSBIN)/goimports
	@echo "generating $*"
	( \
		set -e; \
		cd $(pSrc)/$(call toLowerCase,$*); \
		go generate .; \
		cd ../../../$(pGen); \
		$(TOOLSBIN)/goimports -w $*.go \
	)

# Target is for debugging.
.PHONY: generate-kustomize-builtin-plugins
generate-kustomize-builtin-plugins: $(builtinplugins)

### End kustomize plugin rules.

.PHONY: lint-kustomize
lint-kustomize: install-tools $(builtinplugins)
	cd api; \
	$(TOOLSBIN)/golangci-lint-kustomize -c ../.golangci-kustomize.yml run ./...
	cd kustomize; \
	$(TOOLSBIN)/golangci-lint-kustomize -c ../.golangci-kustomize.yml run ./...
	cd pluginator; \
	$(TOOLSBIN)/golangci-lint-kustomize -c ../.golangci-kustomize.yml run ./...

.PHONY: generate
generate: $(builtinplugins)

.PHONY: test-unit-kustomize-api
test-unit-kustomize-api: $(builtinplugins)
	cd api; go test ./...

.PHONY: test-unit-kustomize-plugins
test-unit-kustomize-plugins:
	./hack/testUnitKustomizePlugins.sh

.PHONY: test-unit-kustomize-cli
test-unit-kustomize-cli:
	cd kustomize; go test ./...

.PHONY: test-unit-kustomize-all
test-unit-kustomize-all: \
	test-unit-kustomize-api \
	test-unit-kustomize-cli \
	test-unit-kustomize-plugins

test-unit-cmd-all:
	./travis/kyaml-pre-commit.sh

test-go-mod:
	./travis/check-go-mod.sh

COVER_FILE=coverage.out

.PHONY: cover
cover:
	# The plugin directory eludes coverage, and is therefore omitted
	cd api && go test ./... -coverprofile=$(COVER_FILE) && \
	go tool cover -html=$(COVER_FILE)

.PHONY:
test-examples-kustomize-against-HEAD: $(MYGOBIN)/kustomize $(TOOLSBIN)/mdrip
	./hack/testExamplesAgainstKustomize.sh HEAD

.PHONY:
test-examples-e2e-kustomize: $(MYGOBIN)/mdrip $(MYGOBIN)/kind
	( \
		set -e; \
		/bin/rm -f $(TOOLSBIN)/kustomize; \
		echo "Installing kustomize from ."; \
		cd kustomize; go install .; cd ..; \
		echo "Installing resource from ."; \
		cd cmd/resource; go install .; cd ../..; \
		./hack/testExamplesE2EAgainstKustomize.sh .; \
	)

.PHONY:
test-examples-kustomize-against-latest: $(TOOLSBIN)/mdrip
	( \
		set -e; \
		/bin/rm -f $(TOOLSBIN)/kustomize; \
		echo "Installing kustomize from latest."; \
		GO111MODULE=on GOBIN=$(TOOLSBIN) go install sigs.k8s.io/kustomize/kustomize/v3; \
		./hack/testExamplesAgainstKustomize.sh latest; \
		echo "Reinstalling kustomize from HEAD."; \
		/bin/rm -f $(TOOLSBIN)/kustomize; \
		cd kustomize; GOBIN=$(MYGOBIN) go install .; \
	)

# linux only.
# This is for testing an example plugin that
# uses kubeval for validation.
# Don't want to add a hard dependence in go.mod file
# to github.com/instrumenta/kubeval.
# Instead, download the binary.
$(TOOLSBIN)/kubeval:
	( \
		set -e; \
		d=$(shell mktemp -d); cd $$d; \
		wget https://github.com/instrumenta/kubeval/releases/latest/download/kubeval-linux-amd64.tar.gz; \
		tar xf kubeval-linux-amd64.tar.gz; \
		mv kubeval $(TOOLSBIN); \
		rm -rf $$d; \
	)

# linux only.
# This is for testing an example plugin that
# uses helm to inflate a chart for subsequent kustomization.
# Don't want to add a hard dependence in go.mod file
# to helm.
# Instead, download the binary.
$(TOOLSBIN)/helm:
	( \
		set -e; \
		d=$(shell mktemp -d); cd $$d; \
		wget https://storage.googleapis.com/kubernetes-helm/helm-v2.16.0-linux-amd64.tar.gz; \
		tar -xvzf helm-v2.16.0-linux-amd64.tar.gz; \
		mv linux-amd64/helm $(TOOLSBIN); \
		rm -rf $$d \
	)

$(TOOLSBIN)/kind:
	( \
        	set -e; \
        	d=$(shell mktemp -d); cd $$d; \
        	wget -O ./kind https://github.com/kubernetes-sigs/kind/releases/download/v0.7.0/kind-$(shell uname)-amd64; \
        	chmod +x ./kind; \
        	mv ./kind $(TOOLSBIN); \
        	rm -rf $$d; \
	)

.PHONY: fmt-api
fmt-api:
	cd api; go fmt ./...

.PHONY: fmt-kustomize
fmt-kustomize:
	cd kustomize; go fmt ./...

.PHONY: fmt-pluginator
fmt-pluginator:
	cd pluginator; go fmt ./...

.PHONY: fmt-plugins
fmt-plugins:
	cd plugin/builtin/prefixsuffixtransformer && go fmt ./...
	cd plugin/builtin/replicacounttransformer && go fmt ./...
	cd plugin/builtin/patchstrategicmergetransformer && go fmt ./...
	cd plugin/builtin/imagetagtransformer && go fmt ./...
	cd plugin/builtin/namespacetransformer && go fmt ./...
	cd plugin/builtin/labeltransformer && go fmt ./...
	cd plugin/builtin/legacyordertransformer && go fmt ./...
	cd plugin/builtin/patchtransformer && go fmt ./...
	cd plugin/builtin/configmapgenerator && go fmt ./...
	cd plugin/builtin/inventorytransformer && go fmt ./...
	cd plugin/builtin/annotationstransformer && go fmt ./...
	cd plugin/builtin/secretgenerator && go fmt ./...
	cd plugin/builtin/patchjson6902transformer && go fmt ./...
	cd plugin/builtin/hashtransformer && go fmt ./...

.PHONY: fmt
fmt: fmt-api fmt-kustomize fmt-pluginator fmt-plugins

.PHONY: modules
modules:
	./hack/doGoMod.sh tidy

## --------------------------------------
## Binaries
## --------------------------------------

.PHONY: build
build:
	cd pluginator && go build -o $(PLUGINATOR_NAME) .
	cd kustomize && go build -o $(KUSTOMIZE_NAME) ./main.go

.PHONY: install
install: generate
	cd pluginator && GOBIN=$(TOOLSBIN) go install $(PWD)/pluginator
	cd kustomize && GOBIN=$(MYGOBIN) go install $(PWD)/kustomize

.PHONY: clean
clean:
	go clean --cache
	rm -f api/$(COVER_FILE)
	rm -f $(builtinplugins)
	rm -f $(MYGOBIN)/kustomize
	rm -f $(TOOLSBIN)/goimports
	rm -f $(TOOLSBIN)/golangci-lint-kustomize
	rm -f $(TOOLSBIN)/golangci-lint
	rm -f $(TOOLSBIN)/mdrip
	rm -f $(TOOLSBIN)/pluginator
	rm -f $(TOOLSBIN)/stringer

.PHONY: nuke
nuke: clean
	sudo rm -rf $(shell go env GOPATH)/pkg/mod/sigs.k8s.io
