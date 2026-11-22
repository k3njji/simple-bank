
-- Create users table
CREATE TABLE "users" (
    "username" varchar PRIMARY KEY,
    "hashed_password" varchar NOT NULL,
    "full_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Link accounts to users
ALTER TABLE "accounts"
ADD CONSTRAINT "accounts_owner_fkey"
FOREIGN KEY ("owner")
REFERENCES "users" ("username")
DEFERRABLE INITIALLY IMMEDIATE;

-- Prevent a user from having multiple accounts in the same currency
ALTER TABLE "accounts"
ADD CONSTRAINT "owner_currency_key"
UNIQUE ("owner", "currency");
