FROM golang:1.10.7-alpine3.7

RUN mkdir -p "$GOPATH/src/github.com/opensteel/authserver"
COPY ./ $GOPATH/src/github.com/opensteel/authserver

RUN cd "$GOPATH/src/github.com/opensteel/authserver/cmd"; go build

ENV DatabaseIp="172.18.238.10"
ENV DatabasePassw="mysecret"
ENV DatabaseName="authdb"
ENV DatabaseUser="postgres"
ENV ServerPort="8080"

WORKDIR "$GOPATH/src/github.com/opensteel/authserver/cmd"
ENTRYPOINT exec ./cmd -db="user=$DatabaseUser password=$DatabasePassw host=$DatabaseIp dbname=$DatabaseName sslmode=disable" -port=":8989"