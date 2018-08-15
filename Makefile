PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation')
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_TAGS = netgo ledger
BBUILD_FLAGS = -tags "${BUILD_TAGS}" -ldflags "-X github.com/icheckteam/ichain/version.GitCommit=${COMMIT_HASH}"
GCC := $(shell command -v gcc 2> /dev/null)
LEDGER_ENABLED ?= true

all: get_vendor_deps install test

########################################
### Build/Install

check-ledger: 
ifeq ($(LEDGER_ENABLED),true)
ifndef GCC
$(error "gcc not installed for ledger support, please install")
endif
else
TMP_BUILD_TAGS := $(BUILD_TAGS)
BUILD_TAGS = $(filter-out ledger, $(TMP_BUILD_TAGS))
endif

########################################
### Build
# This can be unified later, here for easy demos
build: check-ledger
ifeq ($(OS),Windows_NT)
	go build $(BUILD_FLAGS) -o build/ichaind.exe ./cmd/ichaind
	go build $(BUILD_FLAGS) -o build/ichaincli.exe ./cmd/ichaincli
else
	go build $(BUILD_FLAGS) -o build/ichaind ./cmd/ichaind
	go build $(BUILD_FLAGS) -o build/ichaincli ./cmd/ichaincli
endif

build-linux:
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build


install: 
	go install $(BUILD_FLAGS) ./cmd/ichaind
	go install $(BUILD_FLAGS) ./cmd/ichaincli

########################################
### Tools & dependencies

check_tools:
	cd tools && $(MAKE) check_tools

check_dev_tools:
	cd tools && $(MAKE) check_dev_tools

update_tools:
	cd tools && $(MAKE) update_tools

update_dev_tools:
	cd tools && $(MAKE) update_dev_tools

get_tools:
	cd tools && $(MAKE) get_tools

get_dev_tools:
	cd tools && $(MAKE) get_dev_tools

get_vendor_deps:
	@echo "--> Running dep ensure"
	@dep ensure -v


########################################
### Documentation

godocs:
	@echo "--> Wait a few seconds and visit http://localhost:6060/pkg/github.com/icheckteam/ichain/types"
godoc -http=:6060


########################################
### Testing

test: test_unit

test_unit:
	@go test $(PACKAGES_NOSIMULATION)

test_race:
	@go test -race $(PACKAGES_NOSIMULATION)

test_sim:
	@echo "Running individual module simulations."
	@go test $(PACKAGES_SIMTEST) -v
	@echo "Running full Gaia simulation. This may take several minutes."
	@echo "Pass the flag 'SimulationSeed' to run with a constant seed."
	@echo "Pass the flag 'SimulationNumKeys' to run with the specified number of keys."
	@echo "Pass the flag 'SimulationNumBlocks' to run with the specified number of blocks."
	@echo "Pass the flag 'SimulationBlockSize' to run with the specified block size (operations per block)."
	@go test ./cmd/gaia/app -run TestFullGaiaSimulation -SimulationEnabled=true -SimulationBlockSize=200 -v

test_cover:
	@bash tests/test_cover.sh

test_lint:
	gometalinter.v2 --config=tools/gometalinter.json ./...
	!(gometalinter.v2 --disable-all --enable='errcheck' --vendor ./... | grep -v "client/")
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	dep status >> /dev/null
	!(grep -n branch Gopkg.toml)

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs misspell -w

benchmark:
	@go test -bench=. $(PACKAGES_NOSIMULATION)


########################################
### Local validator nodes using docker and docker-compose

build-docker-gaiadnode:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/ichaind/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/ichaind:Z tendermint/ichaindnode testnet --v 4 --o . --starting-ip-address 192.168.10.2 ; fi
	docker-compose up

# Stop testnet
localnet-stop:
	docker-compose down

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build install