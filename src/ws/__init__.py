import os
import logging

from flask import Flask
import click
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy.orm import DeclarativeBase
from sqlalchemy.types import PickleType
import dill
from . import crypto

class Base(DeclarativeBase):
    pass

class DillPickleType(PickleType):
    impl = PickleType.impl

    def process_bind_param(self, value, dialect):
        if value is not None:
            return dill.dumps(value, dill.HIGHEST_PROTOCOL, fmode=dill.CONTENTS_FMODE, recurse=True)
        return None

    def process_result_value(self, value, dialect):
        if value is not None:
            return dill.loads(value)
        return None

db = SQLAlchemy(model_class=Base)

public_key_ring = crypto.get_public_key_ring(os.path.join(os.environ.get('PWD'), 'instance', 'keys'))

def create_app(debug_config=None) -> Flask:
    app = Flask(__name__,
        instance_path=os.path.join(os.environ.get('PWD'),'instance'),
        instance_relative_config=True)

    try:
        os.makedirs(app.instance_path)
    except OSError:
        pass

    app.config.from_mapping(
        SECRET_KEY='dev_secret_key',
        NAMESPACE='default',
        SQLALCHEMY_DATABASE_URI='sqlite:///ws.sqlite',
    )

    app.config.from_pyfile('../config.ini', silent=True)

    if not debug_config is None:
        app.config.from_mapping(debug_config)

    if app.config.get('DEBUG'):
        app.logger.setLevel(logging.DEBUG)
    else:
        loglevel = logging.getLevelName(app.config.get('LOG_LEVEL'))
        if isinstance(loglevel, int):
            app.logger.setLevel(loglevel)
    app.logger.debug(app.config)

    db.init_app(app)
    # app.teardown_appcontext(lambda exception: db.session.remove)

    with app.app_context():
        from .blueprints.grader import user, group, question, answer, module
        app.register_blueprint(user.bp)
        app.register_blueprint(group.bp)
        app.register_blueprint(question.bp)
        app.register_blueprint(module.bp)
        app.register_blueprint(answer.bp)

        from .blueprints.jwt import auth

    app.cli.add_command(init_db_command)

    return app

def init_db():
    db.create_all()

@click.command('init-db')
def init_db_command():
    """Initialize the database."""
    init_db()
    click.echo('Initialized the database.')

