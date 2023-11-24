CREATE TABLE IF NOT EXISTS users (
  id bigserial,
  first_name text NOT NULL ,
  last_name text NOT NULL ,
  email text,
  
  _created_at timestamp DEFAULT now() NOT NULL ,
  _modified_at timestamp DEFAULT now() NOT NULL ,
  PRIMARY KEY(id)
);


CREATE TABLE IF NOT EXISTS auth (
  user_id bigserial,
  login text not NULL,
  password text not NULL
);
