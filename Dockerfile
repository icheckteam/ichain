# Simple usage with a mounted data directory:
# > docker build -t ichain .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.ichaind:/root/.ichaind -v ~/.ichaincli:/root/.ichaincli ichain ichaind init
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.ichaind:/root/.ichaind -v ~/.ichaincli:/root/.ichaincli ichain ichaind start
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES make git libc-dev bash gcc linux-headers eudev-dev

# Set working directory for the build
WORKDIR /go/src/github.com/icheckteam/ichain

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
    make get_tools && \
    make get_vendor_deps && \
    make build && \
    make install

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/ichaind /usr/bin/ichaind
COPY --from=build-env /go/bin/ichaincli /usr/bin/ichaincli

# Run gaiad by default, omit entrypoint to ease using container with ichaincli
CMD ["ichaind"]