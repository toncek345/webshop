CREATE TABLE news (
  id serial NOT NULL PRIMARY KEY,
  header text NOT NULL,
  text text NOT NULL,
  image_key text NOT NULL
);
