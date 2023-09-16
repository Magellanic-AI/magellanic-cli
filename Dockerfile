FROM golang:1.21-alpine3.18 AS builder

WORKDIR $GOPATH/src/magellanic/magellanic-cli/
COPY . .

RUN go get -d -v

RUN go build -o /go/bin/magellanic-cli

FROM scratch

COPY --from=builder /go/bin/magellanic-cli /go/bin/magellanic-cli
ENTRYPOINT ["/go/bin/magellanic-cli"]
