import os
import logging

from flask import Flask

def create_app(debug_config=None) -> Flask:
    app = Flask(__name__, instance_path=os.path.join(os.environ.get('PWD'),'instance'), instance_relative_config=True)
    app.config.from_mapping(
        SECRET_KEY='dev_secret_key',
        NAMESPACE='default',
        DATABASE=os.path.join(app.instance_path, 'ws.sqlite'),
    )

    app.config.from_pyfile('../config.ini', silent=True)

    if not debug_config is None:
        app.config.from_mapping(debug_config)

    try:
        os.makedirs(app.instance_path)
    except OSError:
        pass

    from . import db
    db.init_app(app)

    from . import user
    app.register_blueprint(user.bp)

    from . import submission
    app.register_blueprint(submission.bp)

    from . import question
    app.register_blueprint(question.bp)

    if app.config.get('DEBUG'):
        app.logger.setLevel(logging.DEBUG)
    else:
        loglevel = logging.getLevelName(app.config.get('LOG_LEVEL'))
        if isinstance(loglevel, int):
            app.logger.setLevel(loglevel)
    app.logger.debug(app.config)

    return app
