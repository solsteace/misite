-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE "articles" ADD COLUMN "serie_id" INTEGER;
ALTER TABLE "articles" ADD COLUMN "serie_order" SMALLINT;

UPDATE articles
SET 
    serie_id = article_series.serie_id,
    serie_order = article_series."order"
FROM article_series
WHERE articles.id = article_series.article_id;

ALTER TABLE "articles" ADD FOREIGN KEY ("serie_id") REFERENCES "series"("id");
ALTER TABLE "articles" ALTER COLUMN "serie_id" SET NOT NULL;
ALTER TABLE "articles" ALTER COLUMN "serie_order" SET NOT NULL;
ALTER TABLE "articles" ADD UNIQUE("id", "serie_id");
ALTER TABLE "articles" ADD UNIQUE("serie_id", "serie_order");

DROP TABLE "article_series";

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

CREATE TABLE "article_series"(
    "id" SERIAL PRIMARY KEY,
    "article_id" INTEGER NOT NULL,
    "serie_id" INTEGER NOT NULL,
    "order" SMALLINT NOT NULL,

    FOREIGN KEY ("article_id")
        REFERENCES "articles"("id")
        ON DELETE CASCADE,
    FOREIGN KEY ("serie_id")
        REFERENCES "series"("id")
        ON DELETE CASCADE,

    UNIQUE("serie_id", "article_id"),
    UNIQUE("serie_id", "order"));

INSERT INTO "article_series"("article_id", "serie_id", "order")
SELECT "id", "serie_id", "serie_order"
FROM "articles";

ALTER TABLE "articles" DROP COLUMN "serie_id" CASCADE;
ALTER TABLE "articles" DROP COLUMN "serie_order" CASCADE;