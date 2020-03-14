CREATE TABLE authentication (
  id serial NOT NULL PRIMARY KEY,
  user_id int NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  valid_until timestamp,
  token uuid NOT NULL
);
