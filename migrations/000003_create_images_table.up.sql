CREATE TABLE images (
  id serial NOT NULL PRIMARY KEY,
  product_id bigint NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  name text NOT NULL
);
