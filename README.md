# Pygrader

## Introduction

## First steps

* _flask --app ws init-db_, or
* _flask --app ws reset-db_
* edit ./instance/config.ini
* _uswgi --ini wsgi.ini_
* _uswgi --stop ./instance/pidfile.pid_

Note that _init-db_ will not create existing tables, wherehas _reset-db_ will drop the tables first and recreate them.
