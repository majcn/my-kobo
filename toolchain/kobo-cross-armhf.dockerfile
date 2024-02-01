# source: https://raw.githubusercontent.com/ewpratten/kobo-rs/master/toolchain/kobo-cross-armhf.dockerfile
FROM ghcr.io/cross-rs/arm-unknown-linux-musleabihf:edge

RUN dpkg --add-architecture armhf

# Configure rust to static compile
ENV RUSTFLAGS='-C target-feature=+crt-static'
