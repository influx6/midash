#!/bin/sh
mysql-start 

echo "------------------------------START MIGRATIONS------------------------"

echo "Attempting to use mysql for User: $MYSQL_USER BY Password: $MYSQL_PASSWORD"
cat /migrations/migrations.sql | mysql --user=$MYSQL_USER --password=$MYSQL_PASSWORD

echo "------------------------------START APP------------------------"
midash