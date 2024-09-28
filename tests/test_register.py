import os
import http.client
import json
import base64
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives import hashes, serialization
import dill as dill
import dill.settings

dill.settings['recurse']=True
dill.settings['ignore']=True
dill.settings['fmode']=dill.HANDLE_FMODE


if os.path.exists("key.pem"):
    with open("key.pem", "rb") as f:
        private_key_pem = f.read()
    private_key = serialization.load_pem_private_key(
        private_key_pem,
        password=None  # Use the password if the private key is encrypted
    )
else:
    # Generate an ECDSA private key
    private_key = ec.generate_private_key(ec.SECP256R1())
    private_key_pem = private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.TraditionalOpenSSL,
        encryption_algorithm=serialization.NoEncryption()  # set a password here if needed
    )
    with open('key.pem', 'wb') as f:
        f.write(private_key_pem)

# Generate the public key from the private key
public_key = private_key.public_key()

# Serialize the public key to PEM format
public_key_pem = public_key.public_bytes(
    encoding=serialization.Encoding.PEM,
    format=serialization.PublicFormat.SubjectPublicKeyInfo
)

# Base64 encode the public key
public_key_encoded = base64.b64encode(public_key_pem).decode()

namespace='namespace:test'
# TEST - create user
# Message and nonce
username = "john.doe"
email = "john.doe@rip.org"
nonce = 12345

# Construct the data to be signed
signed_data = f"{username}:{email}:{public_key_encoded}:{nonce}".encode()

# Sign the data
signature = private_key.sign(
    signed_data,
    ec.ECDSA(hashes.SHA256())
)

# Base64 encode the signature
signature_encoded = base64.b64encode(signature).decode()

# Construct the JSON payload
json_body = {
    "username": username,
    "email": email,
    "public_key": public_key_encoded,
    "nonce": nonce,
    "signature": signature_encoded
}

# Convert the dictionary to a JSON string
json_data = json.dumps(json_body)

# Set up the connection
conn = http.client.HTTPConnection("localhost", 8080)

# Send the POST request
conn.request("PUT", f"/v1/user/{namespace}", body=json_data, headers={
    "Content-Type": "application/json"
})

# Get the response
response = conn.getresponse()
data = response.read()
conn.close()

# Print the response
print(response.status, response.reason)
print(data.decode())

# TEST - create group

# Construct the data to be signed
public_key_id = email
groupname = 'alphagroup'
signed_data = f"{public_key_id}:{groupname}:{nonce}".encode()

# Sign the data
signature = private_key.sign(
    signed_data,
    ec.ECDSA(hashes.SHA256())
)

# Base64 encode the signature
signature_encoded = base64.b64encode(signature).decode()

# Construct the JSON payload
json_body = {
    "name": groupname,
    "public_key_id": public_key_id,
    "nonce": nonce,
    "signature": signature_encoded
}

# Convert the dictionary to a JSON string
json_data = json.dumps(json_body)

# Set up the connection
conn = http.client.HTTPConnection("localhost", 8080)

# Send the POST request
conn.request("PUT", f"/v1/group/{namespace}", body=json_data, headers={
    "Content-Type": "application/json"
})

# Get the response
response = conn.getresponse()
data = response.read()
conn.close()

# Print the response
print(response.status, response.reason)
print(data.decode())

# TEST - create module
# Data and nonce
module = 'module_1'
admin = groupname
nonce = 23456
signed_data = f"{public_key_id}:{module}:{admin}:{nonce}".encode()

# Sign the data
signature = private_key.sign(
    signed_data,
    ec.ECDSA(hashes.SHA256())
)

# Base64 encode the signature
signature_encoded = base64.b64encode(signature).decode()

# Construct the JSON payload
json_body = {
    "name": module,
    "admin_id": admin,
    "nonce": nonce,
    "public_key_id": public_key_id,
    "signature": signature_encoded
}

# Convert the dictionary to a JSON string
json_data = json.dumps(json_body)

# Set up the connection
conn = http.client.HTTPConnection("localhost", 8080)

# Send the POST request
conn.request("PUT", f"/v1/module/{namespace}", body=json_data, headers={
    "Content-Type": "application/json"
})

# Get the response
response = conn.getresponse()
data = response.read()
conn.close()

# Print the response
print(response.status, response.reason)
print(data.decode())


# TEST - create question
# Data and nonce
name = 'q_2.1'
method = 'module_1_q_2_1'
max_try=10
nonce = 23456
data = base64.b64encode(dill.dumps(method, dill.DEFAULT_PROTOCOL)).decode()
signed_data = f"{public_key_id}:{module}:{name}:{data}:{max_try}:{nonce}".encode()

# Sign the data
signature = private_key.sign(
    signed_data,
    ec.ECDSA(hashes.SHA256())
)

# Base64 encode the signature
signature_encoded = base64.b64encode(signature).decode()

# Construct the JSON payload
json_body = {
    "max_try": max_try,
    "data": data,
    "nonce": nonce,
    "public_key_id": public_key_id,
    "signature": signature_encoded
}

# Convert the dictionary to a JSON string
json_data = json.dumps(json_body)

# Set up the connection
conn = http.client.HTTPConnection("localhost", 8080)

# Send the POST request
conn.request("PUT", f"/v1/question/{namespace}/{module}/{name}", body=json_data, headers={
    "Content-Type": "application/json"
})

# Get the response
response = conn.getresponse()
data = response.read()
conn.close()

# Print the response
print(response.status, response.reason)
print(data.decode())

# TEST 3 - send an answer (factorial)
def factorial(x):
    if x > 1:
       return x * factorial(x-1)
    else:
       return 1

# Data and nonce
module = 'module_1'
name = 'q_2.1'
nonce = 34567
data = base64.b64encode(dill.dumps(factorial, dill.DEFAULT_PROTOCOL)).decode()
signed_data = f"{public_key_id}:{module}:{name}:{groupname}:{data}:{nonce}".encode()

# Sign the data
signature = private_key.sign(
    signed_data,
    ec.ECDSA(hashes.SHA256())
)

# Base64 encode the signature
signature_encoded = base64.b64encode(signature).decode()

# Construct the JSON payload
json_body = {
    "public_key_id": public_key_id,
    "group_name": groupname,
    "data": data,
    "nonce": nonce,
    "signature": signature_encoded
}

# Convert the dictionary to a JSON string
json_data = json.dumps(json_body)

# Set up the connection
conn = http.client.HTTPConnection("localhost", 8080)

# Send the POST request
conn.request("PUT", f"/v1/answer/{namespace}/{module}/{name}", body=json_data, headers={
    "Content-Type": "application/json"
})

# Get the response
response = conn.getresponse()
data = response.read()

# Print the response
print(response.status, response.reason)
print(data.decode())

# Close the connection
conn.close()
