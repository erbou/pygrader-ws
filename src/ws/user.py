#!/usr/bin/env python3

import functools
import base64
import json
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.backends import default_backend
from flask import Blueprint, g, request, jsonify, current_app, Response
from markupsafe import escape
from ws.db import get_db, ContextManager

bp = Blueprint('user', __name__, url_prefix='/v1/user')

@bp.route('/<namespace>', methods=(['POST']))
def register(namespace):
    try:
        if namespace != current_app.config.get('NAMESPACE'):
            raise RuntimeError(f'Invalid namespace "{namespace}"')

        body = request.get_json()
        current_app.logger.debug(body)
        user = body.get('user')
        key_encoded = body.get('key')
        nonce = body.get('nonce')
        signature_encoded = body.get('signature')
        signed_message = f"{user}:{key_encoded}:{nonce}".encode()
        public_key_bytes = base64.b64decode(key_encoded)
        signature = base64.b64decode(signature_encoded)
        public_key = serialization.load_pem_public_key(public_key_bytes)
        public_key.verify(
            signature,
            signed_message,
            ec.ECDSA(hashes.SHA256())
        )
        db = get_db()
        with ContextManager(db.cursor()) as cr:
            #cr.execute('INSERT OR IGNORE INTO user(username,key) VALUES(?,?) RETURNING id,username,key', (user, public_key_bytes))
            #cr.execute('INSERT OR REPLACE INTO user(username,key) VALUES(?,?) RETURNING id,username,key', (user, public_key_bytes))
            cr.execute('INSERT INTO user(username,key) VALUES(?,?) ON CONFLICT DO UPDATE SET id=id RETURNING id,username,key', (user, public_key_bytes))
            res = cr.fetchone()
            if res:
                db.commit()
                (oid,user,key) = res
                return Response(f'{{ "oid": {oid}, "user": "{user}", "key": "{base64.b64encode(key).decode("utf-8")}" }}', status=200, mimetype='application/json')
            raise RuntimeError("no value")
    except Exception as e:
        current_app.logger.error(f'{e}')

    return Response('{}', status=403, mimetype='application/json')
