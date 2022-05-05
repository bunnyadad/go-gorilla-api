CREATE TABLE IF NOT EXISTS users (
  id SERIAL,
  username varchar(225) NOT NULL,
  password varchar(225) NOT NULL,
  createdat timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updatedat timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  primary key (id)
);

CREATE INDEX index_username ON users (username);
