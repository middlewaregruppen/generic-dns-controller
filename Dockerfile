FROM golang:1.20-alpine AS go-build
COPY . /build
WORKDIR /build/
RUN apk add --no-cache --update git make ca-certificates \
&&  make binary

FROM scratch
LABEL maintainer="@middlewaregruppen (github.com/middlewaregruppen)"
COPY --from=go-build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go-build /build/bin/* /
ENTRYPOINT ["/generic-dns-controller"]
