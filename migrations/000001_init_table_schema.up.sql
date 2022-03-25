CREATE TABLE "orders" (
  "id" varchar PRIMARY KEY,
  "status" int NOT NULL,
  "number" varchar NOT NULL,
  "space_id" varchar NOT NULL,
  "pay_id" varchar,
  "price" float NOT NULL,
  "package_version" varchar NOT NULL,
  "package_price" float NOT NULL,
  "site_code" varchar NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL
);

CREATE TABLE "invoices" (
  "id" varchar PRIMARY KEY,
  "order_id" varchar NOT NULL,
  "status" int NOT NULL,
  "price" float NOT NULL,
  "name" varchar NOT NULL,
  "code" varchar NOT NULL,
  "path" varchar NOT NULL,
  "site_code" varchar NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL
);

CREATE INDEX ON "orders" ("space_id");

CREATE INDEX ON "invoices" ("order_id");

COMMENT ON COLUMN "orders"."id" IS '订单id';

COMMENT ON COLUMN "orders"."status" IS '订单状态';

COMMENT ON COLUMN "orders"."number" IS '订单号';

COMMENT ON COLUMN "orders"."space_id" IS '空间id';

COMMENT ON COLUMN "orders"."pay_id" IS '支付id';

COMMENT ON COLUMN "orders"."price" IS '订单总价';

COMMENT ON COLUMN "orders"."package_version" IS '套餐版本';

COMMENT ON COLUMN "orders"."package_price" IS '套餐价格';

COMMENT ON COLUMN "orders"."site_code" IS '自定义域名标识';

COMMENT ON COLUMN "orders"."created_at" IS '创建时间';

COMMENT ON COLUMN "orders"."updated_at" IS '更新时间';

COMMENT ON COLUMN "invoices"."id" IS '发票id';

COMMENT ON COLUMN "invoices"."order_id" IS '订单id';

COMMENT ON COLUMN "invoices"."status" IS '发票状态';

COMMENT ON COLUMN "invoices"."price" IS '订单总价';

COMMENT ON COLUMN "invoices"."name" IS '发票抬头';

COMMENT ON COLUMN "invoices"."code" IS '发票税号';

COMMENT ON COLUMN "invoices"."path" IS '发票保存路径';

COMMENT ON COLUMN "invoices"."site_code" IS '自定义域名标识';

COMMENT ON COLUMN "invoices"."created_at" IS '创建时间';

COMMENT ON COLUMN "invoices"."updated_at" IS '更新时间';
