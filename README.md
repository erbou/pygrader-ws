# Pygrader

## Introduction

## First steps

* Edit Flask's ./config.ini and uWSGI's ./wsgi.ini
* _flask --app ws init-db_
* _uswgi --ini wsgi.ini_
* _uswgi --stop ./instance/pidfile.pid_

Note that _init-db_ will not create or modify existing tables.

## References

* Flask [configuration](https://flask.palletsprojects.com/en/3.0.x/config/)
* uWSGI [configuration](https://uwsgi-docs.readthedocs.io/en/latest/Configuration.html)
