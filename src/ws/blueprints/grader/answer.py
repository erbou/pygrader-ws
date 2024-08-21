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
from ws import db
from . import models

bp = Blueprint('answer', __name__, url_prefix='/v1/answer')

@bp.route('/<namespace>/<module>/<name>', methods=(['PUT']))
def register(namespace, module, name):
    try:
        if namespace != current_app.config.get('NAMESPACE'):
            raise RuntimeError(f'Invalid namespace "{namespace}"')

        body = request.get_json()
        current_app.logger.debug(body)
        public_key_id = body['public_key_id']
        data_encoded = body['data']
        nonce = body['nonce']
        group_name = body.get('nonce')
        signature_encoded = body['signature']
        data = base64.b64decode(data_encoded)
        signature = base64.b64decode(signature_encoded)

        if group_name is None:
            group_name = public_key_id

        user_instance = db.session.query(models.User).filter_by(username=public_key_id).first()
        if not user_instance:
            raise RuntimeError(f'unknown public_key_id {public_key_id}')

        signed_message = f"{public_key_id}:{module}:{name}:{data_encoded}:{nonce}".encode()
        public_key = serialization.load_pem_public_key(user_instance.public_key)
        public_key.verify(
            signature,
            signed_message,
            ec.ECDSA(hashes.SHA256())
        )

        question_instance = db.session.query(models.Question).filter_by(module=module, name=name).first()
        if not question_instance:
           raise RuntimeError(f'Unknown ref {module}.{name}')

        m = hashlib.sha256()
        m.update(namespace.encode('utf-8'))
        m.update(name.encode('utf-8'))
        m.update(data_encoded.encode('utf-8'))
        digest = m.hexdigest()

        score_instance = db.session.query(models.Score).filter_by(digest=digest).first()
        if not score_instance:
            deserialized = pickle.loads(data)
            if False: ## TODO - validation, make sure it satisfies minimum conditions for scorer namespace.name
                return Response('{}', status=400, mimetype='application/json')
            score_instance = models.Score(digest=digest, data=data)
            db.session.add(score_instance)
            db.session.flush()

        answer_instance = db.session.query(models.Answer).filter_by(user_id=user_instance.id, question_id=question_instance.id, score_id=score_instance.id).first()
        if answer_instance:
            answer_instance.group_name = group_name
        else:
            answer_instance = models.Answer(user_id=user_instance.id, question_id=question_instance.id, score_id=score_instance.id, group_name=group_name)
            db.session.add(answer_instance) 

        try:
            db.session.commit()
        except Exception as e:
            db.session.rollback()
            raise e

        return Response(f'{{ "oid": {answer_instance.id}, "score_id": {score_instance.id}, "digest": "{digest}" }}', status=200, mimetype='application/json')
    except Exception as e:
        current_app.logger.error(f'{e}')

    return Response('{}', status=400, mimetype='application/json')
