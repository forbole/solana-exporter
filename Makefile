COMMIT  := $(shell git log -1 --format='%H')

export GO111MODULE = on

###############################################################################
###                                   All                                   ###
###############################################################################

all: lint test-unit install

###############################################################################
###                                  Build                                  ###
###############################################################################

build: go.sum
	@echo "building solana_exporter binary..."
	@go build -mod=readonly -o build/solana_exporter ./cmd/solana_exporter
.PHONY: build

###############################################################################
###                                 Install                                 ###
###############################################################################

install: go.sum
	@echo "installing solana_exporter binary..."
	@go install -mod=readonly ./cmd/solana_exporter
.PHONY: install

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################
lint:
	golangci-lint run --out-format=tab
.PHONY: lint

lint-fix:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0
.PHONY: lint-fix

clean:
	rm -f tools-stamp ./build/**
.PHONY: clean