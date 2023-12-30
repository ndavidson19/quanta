CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "phone_number" varchar,
  "password_changed_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "last_login_at" timestamptz,
  "login_attempts" int DEFAULT 0,
  "locked_until" timestamptz,
  "reset_token" varchar,
  "reset_token_expires_at" timestamptz
);

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL REFERENCES "users" ("username"),
  "name" varchar NOT NULL,
  "balance" bigint NOT NULL CHECK (balance >= 0),
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "last_updated" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "trades" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL REFERENCES "accounts" ("id"),
  "symbol" varchar NOT NULL,
  "amount" DECIMAL(19,4) NOT NULL CHECK (amount > 0),
  "price" DECIMAL(19,4) NOT NULL,
  "trade_type" varchar NOT NULL,
  "status" varchar NOT NULL DEFAULT 'pending', 
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz DEFAULT (now())
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL REFERENCES "accounts" ("id"),
  "transaction_type" varchar NOT NULL,
  "transaction_amount" DECIMAL(19,4) NOT NULL,
  "transaction_date" timestamptz NOT NULL DEFAULT (now()),
  "related_trade_id" bigint REFERENCES "trades" ("id")
);

CREATE TABLE "financial_instruments" (
  "id" bigserial PRIMARY KEY,
  "symbol" varchar NOT NULL UNIQUE,
  "name" varchar NOT NULL,
  "instrument_type" varchar NOT NULL, -- e.g., 'stock', 'bond', 'future'
  "details" JSON,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "notifications" (
  "id" bigserial PRIMARY KEY,
  "user_id" varchar NOT NULL REFERENCES "users" ("username"),
  "message" text NOT NULL,
  "notification_date" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "deposits" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "portfolio_balances" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "symbol" varchar,
  "amount" int,
  "last_updated" timestamptz DEFAULT (now())
);


CREATE TABLE "audit_logs" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL REFERENCES "accounts" ("id"),
  "action" varchar NOT NULL,
  "action_type" varchar NOT NULL,
  "action_description" varchar NOT NULL,
  "performed_by" varchar NOT NULL REFERENCES "users" ("username"),
  "affected_record" bigint,
  "timestamp" timestamptz DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "deposits" ("account_id");

CREATE INDEX ON "trades" ("account_id");

CREATE INDEX ON "portfolio_balances" ("account_id");

CREATE INDEX ON "trades" ("account_id", "status");

COMMENT ON COLUMN "trades"."amount" IS 'must be positive';

ALTER TABLE "deposits" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "trades" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "trades" ADD CONSTRAINT check_trade_type CHECK (trade_type IN ('buy', 'sell'));

ALTER TABLE "portfolio_balances" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "audit_logs" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");
