#!/bin/sh

if [ -d /db/mysql ]; then
  echo "[i] MySQL directory already present, skipping creation"
else
  echo "[i] MySQL data directory not found, creating initial DBs"

  mysql_install_db --user=root > /dev/null

  if [ "$MYSQL_ROOT_PASSWORD" = "" ]; then
    MYSQL_ROOT_PASSWORD=111111
    echo "[i] MySQL root Password: $MYSQL_ROOT_PASSWORD"
  fi

  MYSQL_DATABASE=${MYSQL_DATABASE:-""}
  MYSQL_USER=${MYSQL_USER:-""}
  MYSQL_PASSWORD=${MYSQL_PASSWORD:-""}

  if [ ! -d "/run/mysqld" ]; then
    mkdir -p /run/mysqld
  fi

  tfile=`mktemp`
  if [ ! -f "$tfile" ]; then
      return 1
  fi

  cat << EOF > $tfile
USE mysql;
FLUSH PRIVILEGES;
EOF

  echo "GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY '$MYSQL_ROOT_PASSWORD' WITH GRANT OPTION;" >> $tfile
  echo "GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' WITH GRANT OPTION;" >> $tfile
  echo "UPDATE user SET password=PASSWORD('') WHERE user='root' AND host='localhost';" >> $tfile

  if [ "$MYSQL_USER" != "" ]; then
    echo "[i] Creating user: $MYSQL_USER with password $MYSQL_PASSWORD"
    echo "CREATE USER '$MYSQL_USER'@'%' IDENTIFIED BY '$MYSQL_PASSWORD';" >> $tfile
    echo "CREATE USER '$MYSQL_USER'@'localhost' IDENTIFIED BY '$MYSQL_PASSWORD';" >> $tfile
    echo "CREATE USER '$MYSQL_USER'@'::1' IDENTIFIED BY '$MYSQL_PASSWORD';" >> $tfile
    # echo "UPDATE user SET password=PASSWORD('$MYSQL_PASSWORD') WHERE user='$MYSQL_USER' AND host='localhost';" >> $tfile
    # echo "UPDATE user SET password=PASSWORD('$MYSQL_PASSWORD') WHERE user='$MYSQL_USER' AND host='%';" >> $tfile
  fi


  if [ "$MYSQL_DATABASE" != "" ]; then
    echo "[i] Creating database: $MYSQL_DATABASE"
    echo "CREATE DATABASE IF NOT EXISTS $MYSQL_DATABASE CHARACTER SET utf8 COLLATE utf8_general_ci;" >> $tfile
    echo "GRANT ALL PRIVILEGES ON $MYSQL_DATABASE.* TO '$MYSQL_USER'@'%' IDENTIFIED BY '$MYSQL_PASSWORD';" >> $tfile
    echo "GRANT ALL PRIVILEGES ON $MYSQL_DATABASE.* TO '$MYSQL_USER'@'localhost' IDENTIFIED BY '$MYSQL_PASSWORD';" >> $tfile
  fi

  if [ "$MYSQL_DAEMONIZE" == "true" ]; then 
    /usr/share/mysql/mysql.server start 
    mysql -uroot < $tfile
  else 
    exec /usr/bin/mysqld --user=root --bootstrap --console --verbose=0 < $tfile 
  fi

  # rm -f $tfile
  echo "Done setting up mysql!"
  exit 0
fi

if [ "$MYSQL_DAEMONIZE" == "true" ]; then 
  /usr/share/mysql/mysql.server start
else 
  exec /usr/bin/mysqld --user=root --console 
fi
