FROM alpine

ARG S6_OVERLAY_VERSION=2.2.0.3
ARG HUEKIT_VERSION=0.1.0

ADD https://github.com/just-containers/s6-overlay/releases/download/v$S6_OVERLAY_VERSION/s6-overlay-amd64-installer /tmp/
RUN chmod +x /tmp/s6-overlay-amd64-installer && /tmp/s6-overlay-amd64-installer /

ADD https://github.com/dj95/huekit/releases/download/v${HUEKIT_VERSION}/huekit_linux_amd64.tar.gz /
RUN tar -xzf /huekit_linux_amd64.tar.gz

VOLUME [ "/huekit_data" ]

ENTRYPOINT ["/init"]
CMD ["/huekit"]