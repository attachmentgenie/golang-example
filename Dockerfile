FROM scratch
ARG TARGETPLATFORM
ENTRYPOINT ["/golang-example", "server"]
COPY $TARGETPLATFORM/golang-example /
