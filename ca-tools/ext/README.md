# ---------------------------------------------------------------------------------------------------------------------------
# CA Chain bundle


## Notes
- In examples below, set the DN names /C, /ST, /O, /OU, /CN and /emailAddress as required.
- The policy of the root CA is configured to impose the intermediate CA to have the same /C, /ST and /O values as the root.
  In addition it requires the server and client csr to specify an emailAddress (requester) and a CN (ip/hostname or client id).
- The provided intermediate CA configuration is set to copy_extensions=on so that it can
  import the alt names of the server certs, but that has security implications. The CA must verify that the requester is not
  including dangerous extensions. If this is an issue, the CA should use separate config files for server and client certs
  and take advantage of -extfile to include server configurations with SAN settings separately.

## About the signing keys

Considerations (circa 2024):

* In all examples we assume that the root CA, intermediate CA, servers, and clients have generated their respective keys (in ./ca/{intr,root})
* Examples for generating RSA and Elliptic Curve (EC) are shown below. EC is preferred because it offers the same security as RSA for less bits.
* Do `openssl ecparam -list_curves` to print the list of available EC curves. In doubt, EC prime256v1 (also secp256r1) is recommended.
* If needed, e.g. for batch csr generation, keys can be generated with -newkey and saved with -keyout, e.g.
  * EC:  openssl req -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -keyout sign.key [...]
  * RSA: openssl req -newkey rsa:4096 -keyout sign.key [...]
* Ideally the CA root key should be encrypted with a password. See genrsa example.

### ECC Key

```
mkdir -m 0700 -p ./private/
openssl ecparam -name prime256v1 -genkey -outform PEM -noout -out ./private/key.pem
chmod 0400 ./private/key.pem
```

### RSA Key (protected by a secret)

```
mkdir -m 0700 -p ./private/
openssl rand -base64 18 -out ./private/secret.txt
~openssl genrsa -passout file:./private/secret.txt -out ./private/key.pem 4096~
openssl genpkey -algorithm rsa -outform PEM  -pkeyopt bits:4096 -pass file:./private/secret.txt
chmod 0400 ./private/key.pem
```

### Warning

The commands shown below make use of openssl ca to manage the certificates.
See the manpage openssl-ca about using this tool in production.

-------
## 1. Create self-signed root CA

The provided ca-root config templates the req DN using ENV variables CA_{DN} (1).
You must set the CA_... environment variables or overload them with `-subj "/C=Country code/ST=State/L=City/O=Org/OU=Unit/CN=Common name"` argument.

```
export $(cat default.env)
```

```
openssl req -config ./conf/openssl-ca-root.cnf -key ./ca/root/private/key.pem -new -x509 -days 9132 -extensions v3_ca -out ./ca/root/certs/ca.cert.pem
chmod 0444 ./ca/root/certs/ca.cert.pem
openssl x509 -noout -text -in ./ca/root/certs/ca.cert.pem
```

Note: (1) it's also possible to use template config, e.g. include {{ variable }} place holders and run openssl as
``openssl req -config <(perl -ne 's/{{([^}]+)}}/$ENV{$1}||$1/eg; print $_' config.cnf)`` (and shoot yourself in the foot).

-------
## 2. Create intermediate CA

Create request:

```
openssl req -config openssl-ca-intr.cnf -key ./ca/intr/private/key.pem -new -out ./ca/root/csr/ca.csr.pem -subj "/emailAddress=Email/C=Country code/ST=State/L=City/O=Org/OU=Unit/CN=Common name delegate"
```

Root CA signs request and issue a certs:

```
openssl ca -config openssl-ca-root.cnf -extensions v3_intermediate_ca -days 730 -notext -md sha256 -in ./ca/root/csr/ca.csr.pem -out ./ca/intr/certs/ca.cert.pem
chmod 0444 ./ca/intr/certs/ca.cert.pem
openssl x509 -noout -text -in ./ca/intr/certs/ca.cert.pem
```

Bundle certificate into a chain (append them into a single file)

```
cat ./ca/intr/certs/ca.cert.pem ./ca/root/certs/ca.cert.pem > ./ca/intr/certs/ca-chain.cert.pem
openssl verify -CAfile ./ca/intr/certs/ca-chain.cert.pem ./ca/intr/certs/ca.cert.pem
openssl verify -CAfile ./ca/root/certs/ca-chain.cert.pem ./ca/intr/certs/ca.cert.pem
```

-------
## Server CA

Create server csr

```
openssl req -new -key ./ca/serv/private/key.pem -out ./ca/intr/csr/serv.csr [-subj '...']
```

CA verifies the request and sign certificates
```
openssl req -in ./ca/intr/csr/serv.csr -text
openssl ca -config ./conf/openssl-ca-intr.cnf -in ./ca/intr/csr/serv.csr -cert ./ca/intr/certs/ca-chain.cert.pem -out server.cert.pem -extensions server_cert
```

## Client CA

Create client csr - client should be provided with methods to generate the private key and the csr on their own with zero knowledge of what they are and how they are used.

```
openssl req -noenc -new -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -keyout client.key -out client.csr -subj "/emailAddress=requester@email/CN=subject@email" 
```

CA verify the client request and signs the certificate (e.g. from whitelist of authorized subjects).
```
openssl req -in client.csr -text
openssl x509 -req -in client.csr -CA ./ca/intr//certs/ca-chain.cert.pem -CAkey ./ca/intr/private/key.pem -out client.cert.pem -CAcreateserial -days 185 -sha256 -extfile client_cert_ext.cnf
```

## Test the certificates

### On the server side

If using Beego, copy the server.cert.pem and the server private key to the ./conf/ folder of the web server and modify the configuration file
Set `TrustCaFile` to the CA chain, and set `EnableMutualHTTPS=true` to force mTLS.

```
HTTPSCertFile = ./conf/server.cert.pem
HTTPSKeyFile = ./conf/server.key.pem
TrustCaFile = ./conf/ca-chain.cert.pem
EnableHTTP = false
EnableMutualHTTPS = true
https = 443
```

If https is terminated on the reverse proxy (e.g. nginx), configure the proxy for mTLS and ensure that the client certificate DN is transferred in the header.
See the [ngnix doc](https://nginx.org/en/docs/http/ngx_http_ssl_module.html).

### On the client side

```
curl -vvv --cert client.cert.pem --key client.key --cacert ./ca/intr/certs/ca-chain.cert.pem 'https://localhost:8443/swagger/swagger.yml'
```

Clients can also concatenate their private keys to their certificates. If using curl, the optional passwords can be appended to the certificate paths, e.g. `./client.cert.pem:password`.

```
cat client.cert.pem client.key > client_key.cert.pem
curl -vvv --cert client.cert.pem --cacert ./ca/intr/certs/ca-chain.cert.pem 'https://localhost:8443/swagger/swagger.yml'
```

