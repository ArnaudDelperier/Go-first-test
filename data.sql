BEGIN;

DROP TABLE IF EXISTS "food_category";
DROP TABLE IF EXISTS "food";
DROP TABLE IF EXISTS "category";


CREATE TABLE "food" (
    "id" int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "name" text NOT NULL,
    "quantity" int NOT NULL
);

CREATE TABLE "category" (
    "id" int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "name" text NOT NULL
);


CREATE TABLE "food_category" (
    "food_id" int NOT NULL REFERENCES "food"("id"),
    "category_id" int NOT NULL REFERENCES "category"("id"),
    PRIMARY KEY ("food_id", "category_id")
);

INSERT INTO "food" ("name", "quantity") VALUES
('carot', 2),
('strawberry', 10),
('tomato', 4);

INSERT INTO "category" ("name") VALUES
('vegetable'),
('fruit');

INSERT INTO "food_category" ("food_id", "category_id") VALUES
(1, 1),
(2, 2),
(3, 1);


COMMIT;