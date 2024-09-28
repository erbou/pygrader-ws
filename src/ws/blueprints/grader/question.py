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

bp = Blueprint('question', __name__, url_prefix='/v1/question')

@bp.route('/<namespace>/<module>/<name>', methods=(['PUT']))
def register(namespace,module,name):
    try:
        if namespace != current_app.config.get('NAMESPACE'):
            raise RuntimeError(f'Invalid namespace "{namespace}"')

        body = request.get_json()
        current_app.logger.debug(body)
        data_encoded = body['data']
        nonce = body['nonce']
        max_try = int(body['max_try'])
        public_key_id = body['public_key_id']
        signature_encoded = body['signature']
        data = base64.b64decode(data_encoded)
        signature = base64.b64decode(signature_encoded)

        user_instance = db.session.query(models.User).filter_by(email=public_key_id).first()
        if not user_instance:
            raise RuntimeError(f'Unknown public_key_id {public_key_id}')

        signed_message = f"{public_key_id}:{module}:{name}:{data_encoded}:{max_try}:{nonce}".encode()
        public_key = serialization.load_pem_public_key(user_instance.public_key)
        public_key.verify(
                signature,
                signed_message,
                ec.ECDSA(hashes.SHA256())
        )


        method_deserialized = pickle.loads(data)
        # TODO: validate method

        module_instance = db.session.query(models.Module).filter_by(name=module).first()
        if not module_instance:
            raise RuntimeError(f'Unknown module {module}')

        if not user_instance in module_instance.admin.users:
            raise RuntimeError(f'Owner of {public_key_id} is not admin of module {module}')

        question_instance = db.session.query(models.Question).filter_by(module_id=module_instance.id, name=name).first()
        if question_instance:
            question_instance.grader = data
            question_instance.max_try = max_try
        else:
            question_instance = models.Question(module=module_instance, name=name, grader=data, max_try=max_try)
            db.session.add(question_instance) 

        try:
            db.session.commit()
        except Exception as e:
            db.session.rollback()
            raise e

        return Response(f'{{ "oid": {question_instance.id}, "module": "{question_instance.module}", "name": "{question_instance.name}" }}', status=200, mimetype='application/json')
    except Exception as e:
        current_app.logger.error(f'{e}')

    return Response('{}', status=400, mimetype='application/json')
