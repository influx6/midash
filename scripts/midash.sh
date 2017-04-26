#!/bin/sh
mysql-start 

echo "------------------------------Starting Application------------------------"
midash-bin
echo "------------------------------Ending Application------------------------"