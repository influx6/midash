FROM influx6/mysql-alpine

MAINTAINER Alexander Ewetumo <trinoxf@gmail.com>

COPY . /go/src/github.com/gu-io/midash

WORKDIR /go/src/github.com/gu-io/midash

# Grab lint tools
RUN go get -u -v github.com/alecthomas/gometalinter

# Install missing lint tools
RUN gometalinter --install

# Run go linters
RUN gometalinter --deadline 4m --errors --vendor ./cmd/midash
RUN gometalinter --deadline 4m --errors --vendor ./pkg/db
RUN gometalinter --deadline 4m --errors --vendor ./pkg/handlers
RUN gometalinter --deadline 4m --errors --vendor ./pkg/internals/models
RUN gometalinter --deadline 4m --errors --vendor ./pkg/internals/utils