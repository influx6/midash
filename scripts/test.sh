#!/bin/sh
mysql-start 

echo "------------------------------START MIGRATIONS------------------------"

echo "Attempting to use mysql with User: $MYSQL_USER"
cat /migrations/migrations.sql | mysql --user=$MYSQL_USER --password=$MYSQL_PASSWORD

echo "----------------------------------------------------------------------"