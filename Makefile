
PACKAGES=$(shell go list ./... | grep -v '/vendor/')
COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_FLAGS = -ldflags "-X github.com/icheckteam/ichain/version.GitCommit=${COMMIT_HASH}"
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

update_tools:
	cd tools && $(MAKE) update_tools

get_tools:
	cd tools && $(MAKE) get_tools

get_vendor_deps:
	@rm -rf vendor/
	@echo "--> Running dep ensure"
	@dep ensure -v


########################################
### Documentation

godocs:
	@echo "--> Wait a few seconds and visit http://localhost:6060/pkg/github.com/icheckteam/ichain/types"
godoc -http=:6060


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

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs misspell -w

benchmark:
	@go test -bench=. $(PACKAGES_NOCLITEST)


########################################
### Local validator nodes using docker and docker-compose

build-docker-ichaindnode:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/ichaind/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/ichaind:Z icheckteam/ichaindnode testnet --v 4 --o . --starting-ip-address 192.168.10.2 ; fi
	docker-compose up



# Stop testnet
localnet-stop:
	docker-compose down


########################################
### Remote validator nodes using terraform and ansible

TESTNET_NAME?=remotenet
SERVERS?=4
BINARY=$(CURDIR)/build/ichaind
remotenet-start:
	@if [ -z "$(DO_API_TOKEN)" ]; then echo "DO_API_TOKEN environment variable not set." ; false ; fi
	@if ! [ -f $(HOME)/.ssh/id_rsa.pub ]; then ssh-keygen ; fi
	@if [ -z "`file $(BINARY) | grep 'ELF 64-bit'`" ]; then echo "Please build a linux binary using 'make build-linux'." ; false ; fi
	cd networks/remote/terraform && terraform init && terraform apply -var DO_API_TOKEN="$(DO_API_TOKEN)" -var SSH_PUBLIC_FILE="$(HOME)/.ssh/id_rsa.pub" -var SSH_PRIVATE_FILE="$(HOME)/.ssh/id_rsa" -var TESTNET_NAME="$(TESTNET_NAME)" -var SERVERS="$(SERVERS)"
	cd networks/remote/ansible && ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i inventory/digital_ocean.py -l "$(TESTNET_NAME)" -e BINARY=$(BINARY) -e TESTNET_NAME="$(TESTNET_NAME)" setup-validators.yml
	cd networks/remote/ansible && ansible-playbook -i inventory/digital_ocean.py -l "$(TESTNET_NAME)" start.yml

remotenet-stop:
	@if [ -z "$(DO_API_TOKEN)" ]; then echo "DO_API_TOKEN environment variable not set." ; false ; fi
	cd networks/remote/terraform && terraform destroy -var DO_API_TOKEN="$(DO_API_TOKEN)" -var SSH_PUBLIC_FILE="$(HOME)/.ssh/id_rsa.pub" -var SSH_PRIVATE_FILE="$(HOME)/.ssh/id_rsa"

remotenet-status:
	cd networks/remote/ansible && ansible-playbook -i inventory/digital_ocean.py -l "$(TESTNET_NAME)" status.yml


# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build install