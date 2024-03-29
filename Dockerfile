FROM golang:1.12.7-buster as builder
LABEL maintainer <athanasios@akostopoulos.com>
WORKDIR /go
COPY . /go
RUN go get -d -v ./...
RUN go build -v -o proxy
# final stage
FROM debian:buster
WORKDIR /go
COPY --from=builder /go/proxy /go/
EXPOSE 53
RUN setcap 'cap_net_bind_service=+ep' proxy
RUN groupadd appuser && useradd -r -u 1001 -g appuser appuser
USER appuser
CMD /go/proxy
