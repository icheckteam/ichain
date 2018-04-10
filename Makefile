PACKAGES=$(shell go list ./... | grep -v '/vendor/')
COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_FLAGS = -ldflags "-X github.com/icheckteam/ichain/version.GitCommit=${COMMIT_HASH}"

all: get_vendor_deps build test

########################################
### Build

# This can be unified later, here for easy demos
build:
	go build $(BUILD_FLAGS) -o build/ichaind ./cmd/ichaind
	go build $(BUILD_FLAGS) -o build/ichaincli ./cmd/ichaincli


install: 
	go install $(BUILD_FLAGS) ./cmd/ichaind
	go install $(BUILD_FLAGS) ./cmd/ichaincli


get_vendor_deps:
	@rm -rf vendor/
	@echo "--> Running dep ensure"
	@dep ensure -v



########################################
### Testing

test: test_unit # test_cli

# Must  be run in each package seperately for the visualization
# Added here for easy reference
# coverage:
#	 go test -coverprofile=c.out && go tool cover -html=c.out

test_unit:
	@go test $(PACKAGES)

test_cover:
	@bash tests/test_cover.sh

benchmark:
	@go test -bench=. $(PACKAGES)