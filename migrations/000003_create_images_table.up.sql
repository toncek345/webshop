CREATE TABLE images (
  id serial NOT NULL PRIMARY KEY,
  key text NOT NULL
);

CREATE TABLE news_images (
  id serial NOT NULL PRIMARY KEY,
  image_id bigint NOT NULL REFERENCES images(id),
  news_id bigint NOT NULL REFERENCES news(id)
);

CREATE TABLE products_images (
  id serial NOT NULL PRIMARY KEY,
  image_id bigint NOT NULL REFERENCES images(id),
  product_id bigint NOT NULL REFERENCES products(id)
);
