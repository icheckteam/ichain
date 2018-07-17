# Local Cluster with Docker Compose

## Requirements

- [Install ichain](../../docs/install.md)
- [Install docker](https://docs.docker.com/engine/installation/)
- [Install docker-compose](https://docs.docker.com/compose/install/)

## Build

Build the `ichaind` binary and the `icheckteam/ichaindnode` docker image.

Note the binary will be mounted into the container so it can be updated without
rebuilding the image.

```
cd $GOPATH/src/github.com/icheckteam/ichain

# Build the linux binary in ./build
make build-linux

# Build icheckteam/ichaindnode image
make build-docker-ichaindnode
```

## Run a testnet

To start a 4 node testnet run:

```
make localnet-start
```

The nodes bind their RPC servers to ports 46657, 46660, 46662, and 46664 on the host.
This file creates a 4-node network using the gaiadnode image.
The nodes of the network expose their P2P and RPC endpoints to the host machine on ports 46656-46657, 46659-46660, 46661-46662, and 46663-46664 respectively.

To update the binary, just rebuild it and restart the nodes:

```
make build-linux
make localnet-stop
make localnet-start
```

## Configuration

The `make localnet-start` creates files for a 4-node testnet in `./build` by calling the `ichaindnode testnet` command.

The `./build` directory is mounted to the `/ichaindnode` mount point to attach the binary and config files to the container.

For instance, to create a single node testnet:

```
cd $GOPATH/src/github.com/icheckteam/ichain

# Clear the build folder
rm -rf ./build

# Build binary
make build-linux

# Create configuration
docker run -v `pwd`/build:/ichaind icheckteam/ichaindnode testnet --o . --v 1

#Run the node
docker run -v `pwd`/build:/ichaind icheckteam/ichaindnode
```

## Logging

Log is saved under the attached volume, in the `ichaind.log` file and written on the screen.

## Special binaries

If you have multiple binaries with different names, you can specify which one to run with the BINARY environment variable. The path of the binary is relative to the attached volume.
