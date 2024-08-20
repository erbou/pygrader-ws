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

bp = Blueprint('submission', __name__, url_prefix='/v1/submission')

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

            cr.execute('SELECT id FROM question WHERE module=? AND name=?', (module, name))
            res = cr.fetchone()
            if not res:
               raise RuntimeError(f'Unknown ref {module}.{name}')
            (qid,) = res

            m = hashlib.sha256()
            m.update(namespace.encode('utf-8'))
            m.update(name.encode('utf-8'))
            m.update(data_encoded.encode('utf-8'))
            digest = m.hexdigest()

            cr.execute('INSERT INTO result(digest,data) VALUES(?,?) ON CONFLICT DO UPDATE SET id=id RETURNING id,score', (digest,data)) 
            res = cr.fetchone()
            (rid,score) = res

            if score is None:
                deserialized = pickle.loads(data)
                ## TODO - validation, make sure it satisfies minimum conditions for scorer namespace.name
                if False:
                    return Response('{}', status=400, mimetype='application/json')

            cr.execute('INSERT INTO submission(uid,qid,rid) VALUES(?,?,?) ON CONFLICT DO UPDATE SET id=id RETURNING id', (uid,qid,rid))
            res = cr.fetchone()
            if not res:
                return Response('{}', status=500, mimetype='application/json')
            (oid,) = res

            db.commit()
            return Response(f'{{ "digest": "{digest}", "oid": {oid}, "rid": {rid} }}', status=200, mimetype='application/json')
    except Exception as e:
        current_app.logger.error(f'{e}')

    return Response('{}', status=400, mimetype='application/json')
