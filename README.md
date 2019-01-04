# authserver

[![Go Report Card](https://goreportcard.com/badge/github.com/opensteel/authserver)](https://goreportcard.com/report/github.com/opensteel/authserver)

Data is storing using PostgreSQL DBMS. To work correctly, the database "authdb" must have 2 tables:
+ access_control (user_type VARCHAR, access_mode INTEGER)
+ users (username, first_name, last_name, user_type, e_mail,  password VARCHAR, deleted BOOLEAN)

For more information see example of test database [authdb.sql](forDatabaseDeploy/authdb.sql)

For fast deploying you can use docker-compose. It start PostgreSQL with test data, and the authorization server. It start server on host port 8989. To download and start enter this commands

```
git clone --branch first_branch https://github.com/opensteel/authserver
cd authserver
docker-compose up -d postgres
docker-compose up go
```

If you want to connect to your own database, you can use only [Dockerfile](Dockerfile) for authorization server and build it:
```
git clone --branch first_branch https://github.com/opensteel/authserver
cd authserver
docker build -t goauth:1 .
```
To run it you should write your darabase connection configuaration and publish port. For examaple:
```
docker run -e DatabaseIp="172.17.0.3" -e  DatabaseUser="postgres"-e DatabasePassw="mysecret" -e DatabaseName="authdb"-e -p 9999:8989 -d goauth:1
```
(will start server on host 9999 port)

If you want to deploy wthout docker enter this commands
```
go get https://github.com/opensteel/authserver
cd $GOPATH/src/github.com/opensteel/authserver/cmd
go build 
```
to run server you, as with docker can use server port with -port, and db settings with -db. For Example:
```
./cmd -db="user=postgres password=mysecret host=127.0.0.1 dbname=authdb sslmode=disable" -port=":8989"
```
(for more setting for connecting database see https://godoc.org/github.com/lib/pq)
