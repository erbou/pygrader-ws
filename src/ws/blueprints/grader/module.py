import base64
import json
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.backends import default_backend
from flask import Blueprint, g, request, jsonify, current_app, Response
from markupsafe import escape
from sqlalchemy import Integer, String
from sqlalchemy.orm import Mapped, mapped_column

bp = Blueprint('module', __name__, url_prefix='/v1/module')

from . import models
from . import crypto
from ws import db

@bp.route('/<namespace>', methods=(['PUT']))
def insert_or_update(namespace):
    try:
        (body, user_instance, public_key_id, name, admin_id) = crypto.get_validated_fields(namespace, request, ['public_key_id','name','admin_id'])

        group_instance = db.session.query(models.Group).filter_by(name=admin_id).first()
        if not group_instance:
            raise RuntimeError(f'unknown group {admin_id}')

        if not user_instance in group_instance.users:
            raise RuntimeError(f'forbidden {public_key_id} not authorized in group {admin_id}')

        module_instance = db.session.query(models.Module).filter_by(name=name).first()
        if module_instance:
            if group_instance.id != module_instance.admin_id:
                raise RuntimeError(f'forbidden group {admin_id} not admin of {module_instance.name}')

            if user_instance not in group_instance.users:
                raise RuntimeError(f'reject {user_instance.email} attempt to update {module_instance.name}')
            module_instance.name = name
        else:
            module_instance = models.Module(name=name)
            module_instance.admin_id = group_instance.name
            group_instance.managed_modules.append(module_instance)
            db.session.add(module_instance)

        try:
            db.session.commit()
        except Exception as e:
            db.session.rollback()
            raise e

        return Response(f'{{ "oid": {group_instance.id} }}', status=200, mimetype='application/json')

    except Exception as e:
        current_app.logger.error(f'{e}')

    return Response('{}', status=403, mimetype='application/json')
