DROP TABLE IF EXISTS tasks_labels, labels, tasks, users;

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE labels (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE tasks (
  id SERIAL PRIMARY KEY,
  opened BIGINT NOT NULL,
  closed BIGINT,
  author_id INTEGER NOT NULL REFERENCES users(id),
  assigned_id INTEGER REFERENCES users(id), 
  title TEXT,
  content TEXT
);

CREATE TABLE tasks_labels (
    task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    label_id INTEGER REFERENCES labels(id) ON DELETE CASCADE
);

insert into users(name) values('Alexa');
insert into labels(name) values('one'), ('two');
