#!/usr/bin/env python3

import functools
import base64
import json
import dill as pickle
import hashlib
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.backends import default_backend
from flask import Blueprint, g, request, jsonify, current_app, Response
from markupsafe import escape
from ws.db import get_db, ContextManager

bp = Blueprint('question', __name__, url_prefix='/v1/question')

@bp.route('/<namespace>/<module>/<name>', methods=(['POST']))
def register(namespace, module, name):
    try:
        if namespace != current_app.config.get('NAMESPACE'):
            raise RuntimeError(f'Invalid namespace "{namespace}"')

        body = request.get_json()
        current_app.logger.debug(body)
        user = body.get('user')
        data_encoded = body.get('data')
        data = base64.b64decode(data_encoded)
        nonce = body.get('nonce')
        signature_encoded = body.get('signature')
        signature = base64.b64decode(signature_encoded)
        signed_message = f"{user}:{module}:{name}:{data_encoded}:{nonce}".encode()

        db = get_db()
        with ContextManager(db.cursor()) as cr:
            cr.execute('SELECT id,key FROM user WHERE username=?', (user,))
            res = cr.fetchone()
            if res:
               pass
            else:
               raise RuntimeError(f'Unknown user {user}')

            (uid, public_key_bytes) = res
            public_key = serialization.load_pem_public_key(public_key_bytes)
            public_key.verify(
                signature,
                signed_message,
                ec.ECDSA(hashes.SHA256())
            )
            deserialized = pickle.loads(data)

            cr.execute('INSERT INTO question(module,name,method) VALUES(?,?,?) ON CONFLICT DO UPDATE SET id=id RETURNING id', (module, name, data))
            res = cr.fetchone()
            if not res:
                return Response('{}', status=500, mimetype='application/json')
            (qid,) = res
            db.commit()
            return Response(f'{{ "module": "{module}", "name": "{name}", "oid": {qid} }}', status=200, mimetype='application/json')
    except Exception as e:
        current_app.logger.error(f'{e}')

    return Response('{}', status=400, mimetype='application/json')
