CREATE TABLE
  IF NOT EXISTS users (
    "id" text,
    "phone" text,
    "password" text,
    "secret" text,
    CONSTRAINT users_pkey PRIMARY KEY ("id")
  );
