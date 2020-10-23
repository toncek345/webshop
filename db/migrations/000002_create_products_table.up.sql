CREATE TABLE products (
  id serial NOT NULL PRIMARY KEY,
  price money NOT NULL,
  name text NOT NULL,
  description text NOT NULL
);
