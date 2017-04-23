FROM influx6/mysql-alpine
MAINTAINER Alexander Ewetumo <trinoxf@gmail.com>

WORKDIR  /

COPY .docker_build/midash /bin/midash
RUN chmod +x /bin/midash

COPY migrations /migrations

COPY scripts/midash.sh /bin/midash-bin
RUN chmod +x /bin/midash-bin

ENV MYSQL_DAEMONIZE true

CMD ["midash-bin"]