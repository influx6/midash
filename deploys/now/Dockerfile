FROM influx6/midash

MAINTAINER Alexander Ewetumo <trinoxf@gmail.com>

# Setup a user, so we run as non-root
# RUN adduser -D myuser
# USER myuser 

# Dont ever do this, far better to load such secret through a 
# env files, but since this will be a plugic repo, then we will allow 
# this here.
ENV API_VERSION v1
ENV MYSQL_USER dbUser
ENV MYSQL_IP "0.0.0.0"
ENV MYSQL_DATABASE midash
ENV MYSQL_PASSWORD !thunder#bot
ENV MYSQL_DAEMONIZE true

# Expose the application port
ENV PORT 7080
EXPOSE 7080

CMD ["midash"]