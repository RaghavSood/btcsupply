# Use the NixOS base image
FROM nixos/nix as builder

# Set up the Nix environment
COPY . /src
WORKDIR /src

RUN nix \
    --extra-experimental-features "nix-command flakes" \
    --option filter-syscalls false \
    build

RUN mkdir -p /tmp/nix-store-closure /tmp/litestream
RUN cp -R $(nix-store -qR result/) /tmp/nix-store-closure

ADD https://github.com/benbjohnson/litestream/releases/download/v0.3.13/litestream-v0.3.13-linux-amd64.tar.gz /tmp/litestream.tar.gz
RUN tar -C /tmp/litestream -xzf /tmp/litestream.tar.gz

FROM alpine:3.6 as alpine

RUN apk add -U --no-cache ca-certificates sqlite bash

WORKDIR /app
RUN mkdir -p /app/bin /etc

# Copy /nix/store
COPY --from=builder /tmp/nix-store-closure /nix/store
COPY --from=builder /src/result /app
COPY --from=builder /tmp/litestream/litestream /app/bin/litestream
COPY --from=builder /src/deployment/bin/run.sh /app/bin/run.sh
COPY --from=builder /src/deployment/etc/litestream.yml /etc/litestream.yml
CMD ["/app/bin/run.sh"]
