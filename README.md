# Pygrader

## Introduction

This repository contains the server (_./pygrader_server_) and client (_./pygrader_client_) libraries of the py grader.

## Install the beego and bee framework

* _./pygrader_server_
    - _go get github.com/beego/beego/v2@latest_
    - _go install github.com/beego/bee/v2@latest_
    - _go mod tidy_

* References:
    - [https://github.com/beego/beego](https://github.com/beego/beego)
    - [https://github.com/beego/bee](https://github.com/beego/bee)

## Prepare the environment

### CA chain and server client certificates

* Prepare the x509 certificates per [instructions](./ca-tools/README.md)

### Server

* Copy x509 certificates and server key into the _./pygrader_webserver/conf/_ folder
    - CA root and intermediate chain _ca-chain.cert.pem_
    - server private key for certificate _serv.prvkey.pem_
    - server certificate _serv.cert.pem_
* Create the beego configuration file _./pygrader_webserver/conf/app.conf_
    - See _app.conf.example_ in same folder.

### Client

* Install the client
    - _cd ./pygrader_client_
    - _python3 -m venv .env_
    - _. .env/bin/activate_
    - _pip install -e ._
* Copy the x509 certificates and client key into the client _./pygrader_client/test_ folder
    - _ca-chain.cert.pem_
    - client private key _client.sign.key.pem_
    - client bundle _client.bundle.pem_, concatenation of
        - _client.cert.pem_
        - _client.prvkey.pem_
    - Optionally a client signing private key key.pem, otherwise the key from the certificate bundle is used

## Running the server

* In the _./pygrader_webserver/_ folder
    - _bee generate routers_
    - _bee run_
