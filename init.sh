#!/bin/sh

[[ -d .env ]] || python3 -m venv .env
. .env/bin/activate
pip install flask uwsgi
