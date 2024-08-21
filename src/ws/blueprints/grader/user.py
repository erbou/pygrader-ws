import functools
import base64
import json
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.backends import default_backend
from flask import Blueprint, g, request, jsonify, current_app, Response
from markupsafe import escape
from sqlalchemy import Integer, String
from sqlalchemy.orm import Mapped, mapped_column

bp = Blueprint('user', __name__, url_prefix='/v1/user')

from . import models
from ws import db

@bp.route('/<namespace>', methods=(['PUT']))
def register(namespace):
    try:
        if namespace != current_app.config.get('NAMESPACE'):
            raise RuntimeError(f'Invalid namespace "{namespace}"')

        body = request.get_json()
        current_app.logger.debug(body)

        username = body['username']
        email = body['email']
        key_encoded = body['public_key']
        nonce = body['nonce']
        signature_encoded = body['signature']

        public_key_bytes = base64.b64decode(key_encoded)

        user_instance = db.session.query(models.User).filter_by(username=username).first()

        if user_instance and public_key_bytes != user_instance.public_key:
            if body['reset_token'] != user_instance.public_key:
                raise RuntimeError(f'unauthorized public_key@{username} reset')

        signed_message = f"{username}:{email}:{key_encoded}:{nonce}".encode()
        signature = base64.b64decode(signature_encoded)
        public_key = serialization.load_pem_public_key(public_key_bytes)
        public_key.verify(
            signature,
            signed_message,
            ec.ECDSA(hashes.SHA256())
        )

        if user_instance:
            user_instance.public_key = public_key_bytes
            user_instance.email = email
        else:
            user_instance = models.User(username=username, public_key=public_key_bytes, email=email)
            db.session.add(user_instance) 

        try:
            db.session.commit()
        except Exception as e:
            db.session.rollback()
            raise e

        return Response(f'{{ "oid": {user_instance.id}, "username": "{user_instance.username}" }}', status=200, mimetype='application/json')

    except Exception as e:
        current_app.logger.error(f'{e}')

    return Response('{}', status=403, mimetype='application/json')
