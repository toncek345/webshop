CREATE TABLE images (
  id serial NOT NULL PRIMARY KEY,
  path text NOT NULL,
  product_id bigint NOT NULL REFERENCES products(id)
);
