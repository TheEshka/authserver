FROM postgres:11

ENV POSTGRES_USER="postgres" \ 
 POSTGRES_PASSWORD="mysecret"
#ENV POSTGRES_DB authdb

ADD ./authdb.sql /docker-entrypoint-initdb.d