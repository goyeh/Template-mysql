#!/bin/bash
# no need to create the database, since this is created by admin and only permissions within this db is set.
# mysql -u$MYSQL_USER -p$MYSQL_PASS < $SQL_FOLDER/0000_SetSchema.sql 2>&1 | grep -v password >> deploy.log
