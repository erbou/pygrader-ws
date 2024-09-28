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
        group_name = body['group_name']
        data_encoded = body['data']
        nonce = body['nonce']
        signature_encoded = body['signature']
        data = base64.b64decode(data_encoded)
        signature = base64.b64decode(signature_encoded)

        user_instance = db.session.query(models.User).filter_by(email=public_key_id).first()
        if user_instance is None:
            raise RuntimeError(f'unknown public_key_id {public_key_id}')

        signed_message = f"{public_key_id}:{module}:{name}:{group_name}:{data_encoded}:{nonce}".encode()
        public_key = serialization.load_pem_public_key(user_instance.public_key)
        public_key.verify(
            signature,
            signed_message,
            ec.ECDSA(hashes.SHA256())
        )

        module_instance = db.session.query(models.Module).filter_by(name=module).first()
        if module_instance is None:
            raise RuntimeError(f'Unknown module {module}')

        question_instance = db.session.query(models.Question).filter_by(module_id=module_instance.id, name=name).first()
        if question_instance is None:
           raise RuntimeError(f'Unknown ref {module}.{name}')

        group_instance = db.session.query(models.Group).filter_by(name=group_name).first()
        if group_instance is None:
           raise RuntimeError(f'Unknown group {group_name}')

        m = hashlib.sha256()
        m.update(namespace.encode('utf-8'))
        m.update(name.encode('utf-8'))
        m.update(data_encoded.encode('utf-8'))
        digest = m.hexdigest()

        score_instance = db.session.query(models.Score).filter_by(digest=digest).first()
        if score_instance is None:
            deserialized = pickle.loads(data)
            if False: ## TODO - validation, make sure deserialized satisfies minimum conditions for scorer namespace.name
                return Response('{}', status=400, mimetype='application/json')
            score_instance = models.Score(question_id=question_instance.id, digest=digest, data=data)
            db.session.add(score_instance)
            db.session.flush()
            answer_instance = None
        else:
            answer_instance = db.session.query(models.Answer).filter_by(sender_id=user_instance.id, score_id=score_instance.id).first()

        if answer_instance is None:
            answer_instance = models.Answer(sender_id=user_instance.id, score_id=score_instance.id, group_id=group_instance.id)
            answer_instance.score = score_instance
            db.session.add(answer_instance) 
            current_app.logger.error(f'Answer => {answer_instance}\nScore => {score_instance}')
            db.session.flush()
        else:
            if answer_instance.group_id != group_instance.id:
                return Response('{{ "error": "you cannot answer a question under a different group" }}', status=403, mimetype='application/json')

        try:
            db.session.commit()
        except Exception as e:
            db.session.rollback()
            raise e

        return Response(f'{{ "oid": {answer_instance.id}, "score_id": {score_instance.id}, "digest": "{digest}" }}', status=200, mimetype='application/json')
    except Exception as e:
        current_app.logger.error(f'{e}')

    return Response('{}', status=400, mimetype='application/json')
