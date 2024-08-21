import os
import logging

from flask import Flask
import click
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy.orm import DeclarativeBase
from sqlalchemy.types import PickleType
import dill

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

def create_app(debug_config=None) -> Flask:
    app = Flask(__name__, instance_path=os.path.join(os.environ.get('PWD'),'instance'), instance_relative_config=True)
    app.config.from_mapping(
        SECRET_KEY='dev_secret_key',
        NAMESPACE='default',
        SQLALCHEMY_DATABASE_URI='sqlite:///ws.sqlite',
    )

    app.config.from_pyfile('../config.ini', silent=True)

    if not debug_config is None:
        app.config.from_mapping(debug_config)

    try:
        os.makedirs(app.instance_path)
    except OSError:
        pass

    from .blueprints.grader import user, answer, question
    app.register_blueprint(user.bp)
    app.register_blueprint(question.bp)
    app.register_blueprint(answer.bp)

    db.init_app(app)

    if app.config.get('DEBUG'):
        app.logger.setLevel(logging.DEBUG)
    else:
        loglevel = logging.getLevelName(app.config.get('LOG_LEVEL'))
        if isinstance(loglevel, int):
            app.logger.setLevel(loglevel)
    app.logger.debug(app.config)

    app.cli.add_command(init_db_command)

    return app

def init_db():
    db.create_all()


@click.command('init-db')
def init_db_command():
    """Create the tables."""
    init_db()
    click.echo('Initialized the database.')

