import os
import base64
import re
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.backends import default_backend

def get_key_fingerprint(key) -> str:
    if hasattr(key, 'private_numbers'):
        public_key = key.public_key()
    else:
        public_key = key
    der_encoded_public_key = public_key.public_bytes(
        encoding=serialization.Encoding.DER,
        format=serialization.PublicFormat.SubjectPublicKeyInfo
    )
    digest = hashes.Hash(hashes.SHA256(), backend=default_backend())
    digest.update(der_encoded_public_key)
    fingerprint = digest.finalize()
    return fingerprint

def get_public_key_ring(path:str) -> dict[str,bytes]:
    keys = {}
    os.makedirs(path, mode=0o700, exist_ok=True)
    if os.stat(path).st_mode & 0o077 > 0:
        raise PermissionError(f'{path} permssions are too open (must be 0o700)')
    for key_pem_file in os.listdir(path):
        if os.path.isfile(key_pem_file):
            key_id=re.fullmatch('key_(.*)_pub.pem', key_pem_file)
            if key_id:
                with open(key_pem_file, "rb") as f:
                    public_key_pem = f.read()
                public_key = serialization.load_pem_public_key(public_key_pem, backend=default_backend())
                keys[get_key_fingerprint(public_key).hex()] = public_key
            
    if len(keys) == 0:
        private_key = ec.generate_private_key(ec.SECP256R1())
        private_key_pem = private_key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.TraditionalOpenSSL,
            encryption_algorithm=serialization.NoEncryption()  # set a password here if needed
        )
        public_key = private_key.public_key()
        public_key_pem = public_key.public_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PublicFormat.SubjectPublicKeyInfo
        )
        fingerprint = get_key_fingerprint(public_key).hex()
        with open(os.path.join(path, f'key_{fingerprint}.pem'), 'wb') as f:
            f.write(private_key_pem)
        with open(os.path.join(path, f'key_{fingerprint}_pub.pem'), 'wb') as f:
            f.write(public_key_pem)
    return keys

def verify_signature(message: bytes, public_key_bytes: bytes, signature : bytes):
    public_key = serialization.load_pem_public_key(public_key_bytes)
    public_key.verify(
        signature,
        message,
        ec.ECDSA(hashes.SHA256())
    )
