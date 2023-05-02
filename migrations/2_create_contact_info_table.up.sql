CREATE TABLE
  IF NOT EXISTS contact_info (
    "user_id" text NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    "phone" text NOT NULL,
    "type" text,
    CONSTRAINT contact_info_pkey PRIMARY KEY ("user_id", "phone")
  );
