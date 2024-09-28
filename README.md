# Pygrader

## Introduction

## Install the beego and bee framework

* https://github.com/beego/beego
* https://github.com/beego/bee

## Prepare the environment

* Prepare the certificates per [instructions](./ca-tools/README.md)
* Copy the CA chain and the server certificates into the ./pygrader_webserver_/conf/ folder
    - CA root and intermediate chain `ca-chain.cert.pem`
    - server private key for certificate `serv.prvkey.pem`
    - server certificate `serv.cert.pem`
* Copy the CA chain and the client certificte into the client ./pygrader_client/test folder
    - ca-chain.cert.pem
    - client private key `client.sign.key.pem`
    - client bundle `client.bundle.pem`, concatenation of
        - client.cert.pem
        - client.prvkey.pem
    - optionally a client signing private key key.pem, otherwise the  certificates key is used
* Write the beego configuration file _./pygrader_webserver/conf/app.conf_
    - See _app.conf.example_ in same folder.

## Running the server

* _bee generate routers_
* _bee run_

