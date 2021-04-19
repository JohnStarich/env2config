FROM golang:1.16-alpine as builder
WORKDIR /src
COPY . /src
RUN go build -o ./env2config ./cmd/env2config

FROM scratch
COPY --from=builder /src/env2config /
