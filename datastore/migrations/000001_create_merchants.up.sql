-- loads a new extension into the current database. There must not be an extension of the same name already loaded.
CREATE EXTENSION pgcrypto;

CREATE TABLE merchants (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  number_of_product INTEGER,
  email VARCHAR(300) NOT NULL UNIQUE,
  business_name VARCHAR(100) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
)