FROM golang:1.19-alpine

RUN apk add build-base

RUN apk --no-cache add make git gcc libtool musl-dev ca-certificates dumb-init

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . ./

RUN make build

CMD ["./blockchain-load-monitoring-service"]