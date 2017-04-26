FROM influx6/mysql-alpine
MAINTAINER Alexander Ewetumo <trinoxf@gmail.com>

WORKDIR  /

COPY .docker_build/midash /bin/midash-bin
RUN chmod +x /bin/midash-bin

COPY scripts/midash.sh /bin/midash
RUN chmod +x /bin/midash

ENV MYSQL_DAEMONIZE true

CMD ["midash"]