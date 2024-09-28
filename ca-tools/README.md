# Certificate management

This tutorial demonstrates how to create a root self-signed certificate, an intermediate certificate authority (CA) chain, and server and client certificates.

It assumes that both the root and intermediate certificates are managed by the same entity. Using a CA chain enhances security, as it enables key rotation at the intermediate level without exposing or compromising the root certificate.

The examples in this tutorial use the Elliptic Curve (EC) algorithm for key generation. EC is favored over RSA due to its ability to provide equivalent security with significantly shorter key lengths, resulting in better performance and efficiency.

This tutorial does not follow industry best practices for certificate management. For a more rigorous (but still imperfect) approach, refer to the examples in ./ext/README.md.

## Root CA creates self signed root certificate and intermediate certificate

The `openssl ec` command used in the examples protects the root private key with a password stored in the preexisting file `ca-root.pass.txt`. However, for better security, it is recommended to keep the password and the private key separate.

To achieve this, omit the `-passout` and `-passin` options to prompt the user for a password at runtime, rather than storing it in a file. Alternatively, during development when password protection is not required, you can skip the `openssl ec` command entirely to keep the key unencrypted (in clear text).

* CA Root Self signed certificate (25y)

```
mkdir -m 0700 -p .
openssl ecparam -name prime256v1 -genkey -outform PEM -noout -out ./ca-root.prvkey.pem
openssl ec -in ./ca-root.prvkey.pem -aes256 -passout file:ca-root.pass.txt -out ./ca-root.prvkey.pem
chmod 0400 ./ca-root.prvkey.pem
openssl req -config ./conf/openssl-root.cnf -key ca-root.prvkey.pem -passin file:ca-root.pass.txt -new -x509 -days 9132 -extensions v3_ca -out ./ca-root.cert.pem
openssl x509 -in ./ca-root.cert.pem -noout -text
```

* CA Intermediate certificate

```
mkdir -m 0700 -p .
openssl req -config ./conf/openssl-intr.cnf -days 730 -CA ./ca-root.cert.pem -CAkey ./ca-root.prvkey.pem -passin file:./ca-root.pass.txt -new -x509 -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -noenc -keyout ca-intr.prvkey.pem -out ca-intr.cert.pem
chmod 0400 ./ca-intr.prvkey.pem
openssl x509 -in ./ca-intr.cert.pem -noout -text
```

* Bundle the certificatesS

```
cat ca-intr.cert.pem ca-root.cert.pem > ca-chain.cert.pem
openssl verify -CAfile ./ca-chain.cert.pem ./ca-root.cert.pem
openssl verify -CAfile ./ca-chain.cert.pem ./ca-intr.cert.pem
```

## Create server CA for TLS server Auth

* Server admin create certificate request CSR

```
openssl req -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -noenc -keyout serv.prvkey.pem -out ./serv.csr -subj '/emailAddress=server@email.ch/C=CH/ST=VD/L=Lausanne/O=EPFL/OU=SDSC/CN=127.0.0.1' -addext 'subjectAltName=IP.1:1.2.3.4,DNS.1:localhost' -addext 'keyUsage=critical,digitalSignature,keyEncipherment' -addext 'extendedKeyUsage=critical,serverAuth'
```

* Intermediate CA sign the certificate request

Warning: The example sets `-copy_extensions` copy so that alt names extensions can be imported from the server cert request.
That has security implications, the CA must verify that the csr is not asking for extensions it shouldn't have.

```
openssl req -in ./serv.csr -noout -text
openssl req -config ./conf/openssl-serv.cnf -x509 -in ./serv.csr -CA ./ca-intr.cert.pem -CAkey ./ca-intr.prvkey.pem -CAcreateserial -out serv.cert.pem -copy_extensions copy
openssl x509 -in ./serv.cert.pem -noout -text
```

## Create client CA for mTLS client Auth

* Client creates a certificate request CSR

```
openssl req -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -noenc -keyout client.prvkey.pem -out ./client.csr -subj '/emailAddress=requester@email.ch/C=CH/ST=VD/L=Lausanne/O=EPFL/OU=SDSC/CN=client@email.ch' -addext 'keyUsage=critical,digitalSignature,keyEncipherment' -addext 'extendedKeyUsage=critical,clientAuth'
```

* Intermediate CA sign the certificate request

Warning: The example sets `-copy_extensions` copy so that alt names extensions can be imported from the client cert request.

```
openssl req -in ./client.csr -noout -text
openssl x509 -extfile ./conf/openssl-client.cnf -req -in ./client.csr -CA ./ca-intr.cert.pem -CAkey ./ca-intr.prvkey.pem -out client.cert.pem -ext keyUsage,extendedKeyUsage -copy_extensions copy
openssl x509 -in ./client.cert.pem -noout -text
```

