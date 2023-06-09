CREATE TABLE
  IF NOT EXISTS users (
    "id" text,
    "full_name" text,
    "department" text,
    "position" text,
    "password" text,
    "otp_secret" text,
    CONSTRAINT users_pkey PRIMARY KEY ("id")
  );
