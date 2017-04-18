FROM influx6/mysql-alpine-setup
MAINTAINER Alexander Ewetumo <trinoxf@gmail.com>

WORKDIR  /

COPY .docker_build/midash /bin/midash
RUN chmod +x /bin/midash