import base64
import json
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.backends import default_backend
from flask import current_app
from ws import db, crypto
from . import models

def get_validated_fields(namespace, request, fields, newkey=False):
    if namespace != current_app.config.get('NAMESPACE'):
        raise RuntimeError(f'Invalid namespace "{namespace}"')

    body = request.get_json()
    #values = [f'{body["payload"][f]}' for f in fields]
    #signed_message = ':'.join(values+[body['signature']['nonce']]).encode()
    values = [str(body.get(f)) for f in fields]
    signed_message = ':'.join(values+[str(body['nonce'])]).encode()

    #key_id = body['signature'].get('public_key_id')
    key_id = body.get('public_key_id')
    if key_id is not None:
        user_instance = db.session.query(models.User).filter_by(email=key_id).first()
        if not user_instance:
            raise RuntimeError(f'Unknown public key id: {public_key_id}')
        public_key_bytes = user_instance.public_key
    elif newkey:
        user_instance = None
        public_key_bytes = base64.b64decode(body['public_key'])
    else:
        raise RuntimeError(f'Forbidden')

    crypto.verify_signature(signed_message, public_key_bytes, base64.b64decode(body['signature']))

    return (body, user_instance, *values)
