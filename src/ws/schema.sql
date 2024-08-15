CREATE TABLE IF NOT EXISTS user (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username   TEXT UNIQUE NOT NULL,
  key        BLOB NOT NULL,
  created    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS result (
  id INTEGER  PRIMARY KEY AUTOINCREMENT,
  digest      TEXT NOT NULL,
  score       REAL NOT NULL,
  result      BLOB NOT NULL,
  created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS question (
  id INTEGER  PRIMARY KEY AUTOINCREMENT,
  module      TEXT NOT NULL,
  name        TEXT NOT NULL,
  method      BLOB NOT NULL,
  created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (module,name)
);

CREATE TABLE IF NOT EXISTS submission (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  uid         INTEGER NOT NULL,
  qid         INTEGER NOT NULL,
  rid         INTEGER,
  digest      TEXT NOT NULL,
  data        BLOB NOT NULL,
  created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (uid) REFERENCES user (id),
  FOREIGN KEY (qid) REFERENCES question (id),
  FOREIGN KEY (rid) REFERENCES result (id),
  UNIQUE (uid,qid,digest)
);

CREATE TABLE IF NOT EXISTS scoreboard (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  sid         INTEGER NOT NULL,
  ranking     INTEGER,
  data        BLOB NOT NULL,
  created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (sid) REFERENCES submission (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS result_digest_idx ON result(digest);
CREATE UNIQUE INDEX IF NOT EXISTS submission_digest_idx ON submission(digest);
CREATE UNIQUE INDEX IF NOT EXISTS question_name_idx ON question(module,name);
