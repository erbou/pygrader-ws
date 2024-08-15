import sqlite3
import click
from flask import current_app, g

class ContextManager:
    def __init__(self, obj):
        self.obj = obj

    def __enter__(self):
        return self.obj

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.obj.close()

def get_db():
    if 'db' not in g:
        g.db = sqlite3.connect(
            current_app.config['DATABASE'],
            detect_types=sqlite3.PARSE_DECLTYPES
        )
        g.db.execute('PRAGMA journal_mode=WAL;')
        g.db.row_factory = sqlite3.Row

    return g.db

def close_db(e=None):
    db = g.pop('db', None)
    if db is not None:
        db.close()

def init_db():
    db = get_db()
    with current_app.open_resource('schema.sql') as f:
        db.executescript(f.read().decode('utf8'))

def reset_db():
    db = get_db()
    with current_app.open_resource('reset.sql') as f:
        db.executescript(f.read().decode('utf8'))

def init_app(app):
    app.teardown_appcontext(close_db)
    app.cli.add_command(init_db_command)
    app.cli.add_command(reset_db_command)

@click.command('init-db')
def init_db_command():
    """Create tables."""
    init_db()
    click.echo('Initialized the database.')

@click.command('reset-db')
def reset_db_command():
    """Clear the existing data and create new tables."""
    reset_db()
    init_db()
    click.echo('Reinitialized the database.')
