FROM golang:alpine
MAINTAINER Ewetumo Alexander <trinoxf@gmail.com>

# Create the db directory and expose it for usage.
RUN mkdir /db
VOLUME /db

RUN apk add --update openrc mysql mysql-client && rm -f /var/cache/apk/*

# Copy mysql configuration to the etc folder.
COPY mysql.cnf /etc/mysql/my.cnf

# Copy the mysql startup so we can run up the server.
COPY mysql.sh /bin/mysql-start
RUN chmod +x /bin/mysql-start

# Expose port for mysql usage.
EXPOSE 5880