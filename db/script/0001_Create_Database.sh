#!/bin/bash

mysql -u$MYSQL_USER -p$MYSQL_PASS < $SQL_FOLDER/0001_news.sql 2>&1 | grep -v password >> deploy.log
