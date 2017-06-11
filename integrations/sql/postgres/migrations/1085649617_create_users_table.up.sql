CREATE SEQUENCE user_id_seq;

CREATE TABLE users (
  user_id integer NOT NULL DEFAULT nextval('user_id_seq'),
  name    varchar(40),
  email   varchar(40)
);

ALTER SEQUENCE user_id_seq OWNED BY users.user_id;

