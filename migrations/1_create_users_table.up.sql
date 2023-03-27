CREATE TABLE
  IF NOT EXISTS users (
    "id" text,
    "phone" text NOT NULL UNIQUE,
    "password" text,
    -----------------------
    "otp_enabled" boolean,
    "otp_secret" text,
    -----------------------
    "telegram_user_id" text,
    CONSTRAINT users_pkey PRIMARY KEY ("id")
  );
