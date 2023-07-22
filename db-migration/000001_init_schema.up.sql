CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "deposits" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "trades" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "symbol" varchar,
  "amount" int NOT NULL,
  "price" decimal NOT NULL,
  "trade_type" varchar,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "portfolio_balances" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "symbol" varchar,
  "amount" int,
  "last_updated" timestamptz DEFAULT (now())
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "full_name" varchar,
  "contact_info" varchar,
  "account_id" bigint NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "audit_logs" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "action" varchar,
  "timestamp" timestamptz DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "deposits" ("account_id");

CREATE INDEX ON "trades" ("account_id");

CREATE INDEX ON "portfolio_balances" ("account_id");

COMMENT ON COLUMN "trades"."amount" IS 'must be positive';

ALTER TABLE "deposits" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "trades" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "portfolio_balances" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "audit_logs" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");