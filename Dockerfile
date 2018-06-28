FROM golang:1.10 AS builder
WORKDIR $GOPATH/src/dummy-bc-platform
COPY . .
RUN make build

FROM scratch
COPY --from=builder /go/src/dummy-bc-platform/build/dummy-bc-platform.linux  ./dummy-bc-platform
COPY .env /
ENTRYPOINT ["./dummy-bc-platform"]