import jwt
import base64
import json
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.backends import default_backend
from flask import Blueprint, g, request, jsonify, current_app, Response
from markupsafe import escape
from sqlalchemy import Integer, String
from sqlalchemy.orm import Mapped, mapped_column
from ws import db

bp = Blueprint('auth', __name__, url_prefix='/v1/auth')

@bp.route('/', methods=(['GET']))
def auth():
    auth_header = request.header.get('Authorization')
    if auth_header and auth_header.startswith('Bearer '):
        token = auth_header[len('Bearer '):]
    else:
        abort(401)

    try:
        kid = jwt.get_unverified_header(token).get('kid')
        if not kid:
            raise ValueError('JWT does not contain a kid')
        public_key_bytes = public_keys.get(kid)
        if not public_key_bytes:
            raise ValueError('Unkown public key id {kid}')
        token_decoded = jwt.decode(token, algorithms=['ES256'])
        return Response('', 200)
    except:
        abort(400)

