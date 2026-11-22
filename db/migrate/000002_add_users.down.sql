ALTER TABLE ID EXISTS "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";

ALTER TABLE ID EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_currency_fkey";

DROP TABLE IF EXISTS "users";