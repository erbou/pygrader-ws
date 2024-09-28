# Data models

The data models are implemented using Flask-SQLAlchemy ORM and are, by default, stored in an SQLite3 database.
You can change this default behavior by setting the _SQLALCHEMY_DATABASE_URI_ configuration variable to the desired database URL.
For more information, refer to the Flask-SQLAlchemy [documentation](https://flask-sqlalchemy.palletsprojects.com/en/3.1.x/quickstart/#installation).

A brief description of the models is provided below.

## Models

### Base abstract model

Declares the fields that are common to all the models:

- _id_: unique integer primary key
- _created_: UTC datetime of creation
- _updated_: UTC datetime of last update

### User

- _username_: user identifier (unique, immutable)
- _email_: user email (unique)
- _public_key_: for signatures verification (unique)

Note:
- A secret reset code must be provided in order to modify _public_key_

### Question

- _module_: test identifier (immutable)
- _name_: question identifier within the module (immutable)
- _max_try_: maximum number of distinct answers allowed for a group

Notes:
- The pair (_module_,_name_) must be unique

### Answer

- _user_id_: user who submitted the answer
- _score_id_: foreign key of scoring for the answer (immutable)
- _group_name_: group under which answer is submitted (immutable)

- (_group_name_,_score_id_) is unique - a single record is generated if a group submit the same answer multiple times.

### Score

- _question_id_: foreign key of the question being answered (immutable)
- _data_: a serialized (pickle, or dill) representation of the answer
- _digest_: a hash of the (question,data) pair
- _score_: the grader's score


