import os
import json
import time
import base64
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.backends import default_backend
from datetime import datetime, timezone, tzinfo
from json import JSONEncoder

class DateTimeEncoder(JSONEncoder):
    def default(self, obj):
        if isinstance(obj, datetime):
            return obj.isoformat()
        elif hasattr(obj, '__dict__'):
            return vars(obj)
        return super().default(obj)

class HashField:
    def __init__(self, i):
        self.i = i

class HashableMeta(type):
    def __new__(cls, name, bases, dct):
        # Automatically collect attributes with `HashField` or similar
        clsobj = super().__new__(cls, name, bases, dct)
        clsobj._hash_fields = sorted(
            ((v.i, k) for k, v in dct.items() if isinstance(v, HashField)),
            key=lambda x: x[0]
        )
        return clsobj

class Hashable(metaclass=HashableMeta):
    def hashObject(self):
        return self._hashObject().hex()

    def _hashObject(self):
        digest = hashes.Hash(hashes.SHA256(), backend=default_backend())

        for _, attr_name in self._hash_fields:
            attr_value = getattr(self, attr_name)
            if not attr_value is None:
                digest.update(self._hashValue(attr_value))
        return digest.finalize()

    def _hashValue(self, value):
        if isinstance(value, Hashable):
            return value._hashObject()
        elif isinstance(value, list):
            list_digest = hashes.Hash(hashes.SHA256(), backend=default_backend())
            for item in value:
                list_digest.update(self._hashValue(item))
            return list_digest.finalize()
        elif isinstance(value, int):
            return str(value).encode('utf-8')
        elif isinstance(value, str):
            return value.encode('utf-8')
        elif isinstance(value, float):
            return str(f'{value:.6f}').encode('utf-8')
        elif isinstance(value, datetime):
            return str(int(value.timestamp())).encode('utf-8')
        else:
            raise TypeError(f"Unsupported type for hashing: {value} ({type(value)})")

class Signer():

    def __init__(self, key_path = None, password = None):
        self.k_path = key_path
        self.k_pass = password
        self.prv_key = None
        self.pub_key = None
        self.pub_key_der = None
        self.pub_key_b64 = None
        self.k_fp = None
        self.k_fp_bytes = 20

    @property
    def key_fingerprint_len(self):
        return self.k_fp_bytes

    @key_fingerprint_len.setter
    def key_fingerprint_len(self, nc):
        self.k_fp_bytes = nc

    @property
    def key_path(self):
        return self.k_path

    @key_path.setter
    def key_path(self, path):
        self.k_path = path
        self.k_pass = None
        self.prv_key = None
        self.pub_key = None
        self.pub_key_der = None
        self.pub_key_b64 = None
        self.k_fp = None

    @property
    def private_key(self):
        if self.prv_key is None:
            self.load_key(self.k_path, self.k_pass)
        return self.prv_key

    @property
    def public_key(self):
        if self.pub_key is None:
            self.pub_key = self.private_key.public_key()
        return self.pub_key

    @property
    def public_key_der(self):
        if self.pub_key_der is None:
           self.pub_key_der = self.public_key.public_bytes(
                encoding=serialization.Encoding.DER,
                format=serialization.PublicFormat.SubjectPublicKeyInfo
            )
        return self.pub_key_der

    @property
    def public_key_b64(self):
        if self.pub_key_b64 is None:
            self.pub_key_b64 = base64.b64encode(self.public_key_der)
        return self.pub_key_b64

    @property
    def key_fingerprint(self):
        if self.k_fp is None:
            digest = hashes.Hash(hashes.SHA256(), backend=default_backend())
            digest.update(self.public_key_der)
            self.k_fp = digest.finalize()
        return self.k_fp[:self.k_fp_bytes]

    def signObject(self, obj, params = [], withDigest=False) -> str:
        # Sign the data
        digest_msg = hashes.Hash(hashes.SHA256(), backend=default_backend())
        nonce = f'{time.time():.0f}'
        message = ':'.join([ str(s) for s in params ] + [nonce, obj.hashObject()])
        #print(message)
        digest_msg.update(message.encode())
        digest = digest_msg.finalize()

        signature = self.private_key.sign(
            message.encode(),
            ec.ECDSA(hashes.SHA256())
        )
        
        # Base64 encode the signature
        signature_encoded = base64.b64encode(signature).decode()
        
        # Construct the JSON payload
        json_body = {
            "Payload": obj , #obj.toJson(),
            "Signature": {
                "nonce": nonce,
                "kid": self.key_fingerprint.hex(),
                "ecdas": signature_encoded
            }
        }

        if withDigest:
            json_body['Signature']['digest'] = digest.hex()
        
        # Convert the dictionary to a JSON string
        print(json_body)
        json_data = json.dumps(json_body, cls=DateTimeEncoder)
        return json_data

    def load_key(self, path, password = None, create=True):
        if not os.path.exists(path):
            if create:
                self.gen_key(path, password)
        with open(path, "rb") as f:
            key_pem = f.read()
        key = serialization.load_pem_private_key(
            key_pem,
            password=password
        )
        self.prv_key = key
        self.pub_key = key.public_key()

    def gen_key(self, path, password = None):
        private_key = ec.generate_private_key(ec.SECP256R1())

        if password != None:
            encryption_algorithm = serialization.BestAvailableEncryption(password)
        else:
            encryption_algorithm = serialization.NoEncryption()

        private_key_pem = private_key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.TraditionalOpenSSL,
            encryption_algorithm=encryption_algorithm
        )

        with open(path, 'xb') as f:
            f.write(private_key_pem)

