FROM golang:1.10.7-alpine3.7 as build

WORKDIR /go/src/github.com/opensteel/authserver
COPY . .

RUN cd "$GOPATH/src/github.com/opensteel/authserver/cmd"; go build -o myapp

FROM alpine:latest as runner
WORKDIR /root/

COPY --from=build /go/src/github.com/opensteel/authserver/cmd/myapp . 

ENV DATABASE_IP="127.0.0.1" \
 DATABASE_PASSW="mysecret" \
 DATABASE_NAME="authdb" \
 DATABASE_USER="postgres"

CMD exec ./myapp -db="user=$DATABASE_USER password=$DATABASE_PASSW host=$DATABASE_IP dbname=$DATABASE_NAME sslmode=disable" -port=":8989"