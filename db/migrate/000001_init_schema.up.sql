
-- Create accounts table
CREATE TABLE "accounts" (
    "id" bigserial PRIMARY KEY,
    "owner" varchar NOT NULL,
    "balance" bigint NOT NULL DEFAULT 0,
    "currency" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Create entries table
CREATE TABLE "entries" (
    "id" bigserial PRIMARY KEY,
    "account_id" bigint NOT NULL,
    "amount" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Create transfers table
CREATE TABLE "transfers" (
    "id" bigserial PRIMARY KEY,
    "from_account_id" bigint NOT NULL,
    "to_account_id" bigint NOT NULL,
    "amount" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Create indexes
CREATE INDEX ON "accounts" ("owner");
CREATE INDEX ON "entries" ("account_id");
CREATE INDEX ON "transfers" ("from_account_id");
CREATE INDEX ON "transfers" ("to_account_id");
CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

-- Add foreign keys
ALTER TABLE "entries"
ADD FOREIGN KEY ("account_id")
REFERENCES "accounts" ("id")
DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "transfers"
ADD FOREIGN KEY ("from_account_id")
REFERENCES "accounts" ("id")
DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "transfers"
ADD FOREIGN KEY ("to_account_id")
REFERENCES "accounts" ("id")
DEFERRABLE INITIALLY IMMEDIATE;
