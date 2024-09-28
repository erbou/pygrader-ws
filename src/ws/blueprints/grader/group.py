import base64
import json
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.backends import default_backend
from flask import Blueprint, g, request, jsonify, current_app, Response
from markupsafe import escape
from sqlalchemy import Integer, String
from sqlalchemy.orm import Mapped, mapped_column

bp = Blueprint('group', __name__, url_prefix='/v1/group')

from . import models
from ws import db

@bp.route('/<namespace>', methods=(['PUT']))
def insert_or_update(namespace):
    try:
        if namespace != current_app.config.get('NAMESPACE'):
            raise RuntimeError(f'Invalid namespace "{namespace}"')

        client_cert = request.headers.get('X-Client-Cert')
        if client_cert:
            client_cert_obj = x509.load_pem_x509_certificate(client_cert.encode(), default_backend())
            public_key = client_cert_obj.public_key()
            # TODO: get scope from certificate
            scope=''
        else:
            scope=''

        body = request.get_json()
        current_app.logger.debug(body)

        name = body['name']
        public_key_id = body['public_key_id']
        nonce = body['nonce']
        signature_encoded = body['signature']

        user_instance = db.session.query(models.User).filter_by(email=public_key_id).first()
        if not user_instance:
            raise RuntimeError(f'unknown public_key_id {public_key_id}')

        signed_message = f"{public_key_id}:{name}:{nonce}".encode()
        signature = base64.b64decode(signature_encoded)
        public_key = serialization.load_pem_public_key(user_instance.public_key)
        public_key.verify(
                signature,
                signed_message,
                ec.ECDSA(hashes.SHA256())
        )

        group_instance = db.session.query(models.Group).filter_by(name=name).first()

        if group_instance:
            if user_instance not in group_instance.users:
                raise RuntimeError(f'reject {user_instance.email} attempt to update {group_instance.id}:{user_instance.username}')
            group_instance.name = name
        else:
            group_instance = models.Group(name=name,scope=scope)
            group_instance.users.append(user_instance)
            db.session.add(group_instance)

        try:
            db.session.commit()
        except Exception as e:
            db.session.rollback()
            raise e

        return Response(f'{{ "oid": {group_instance.id} }}', status=200, mimetype='application/json')

    except Exception as e:
        current_app.logger.error(f'{e}')

    return Response('{}', status=403, mimetype='application/json')
