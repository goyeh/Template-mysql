#!/bin/bash
echo "[P1] Version Number is :$1 "
echo "[P2] Database target is :$2 "

echo "Run this script from this script folder"
echo "PWD:$PWD"

source ~/.env

export MYSQL_USER=$MYSQL_USER
export MYSQL_PASS=$MYSQL_PASS
export MYSQL_DB=$MYSQL_DB
export SQL_FOLDER=$PWD/sql
export DBHOST=localhost
export SCRIPTDIR=$PWD/script

for entry in "ls -v $SCRIPTDIR"; do
    echo "Processing : $SCRIPTDIR/$entry "
    $SCRIPTDIR/$entry
done
