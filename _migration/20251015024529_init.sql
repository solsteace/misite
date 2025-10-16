-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE "articles"(
    "id" SERIAL PRIMARY KEY,
    "title" VARCHAR(128) NOT NULL,
    "subtitle" VARCHAR(256) NOT NULL,
    "content" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);

CREATE TABLE "tags"(
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(64) UNIQUE NOT NULL);

CREATE TABLE "series"(
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(128) UNIQUE NOT NULL);

CREATE TABLE "revisions"(
    "id" SERIAL PRIMARY KEY,
    "article_id" INTEGER NOT NULL,
    "message" VARCHAR(256) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY("article_id")
        REFERENCES "articles"("id")
        ON DELETE CASCADE);

CREATE TABLE "article_tags"(
    "id" SERIAL PRIMARY KEY,
    "article_id" INTEGER NOT NULL,
    "tag_id" INTEGER NOT NULL,

    FOREIGN KEY ("article_id")
        REFERENCES "articles"("id")
        ON DELETE CASCADE,
    FOREIGN KEY ("tag_id")
        REFERENCES "tags"("id")
        ON DELETE CASCADE,

    UNIQUE("article_id", "tag_id"));

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

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd


DROP TABLE "article_series";
DROP TABLE "article_tags";
DROP TABLE "revisions";
DROP TABLE "series"
DROP TABLE "tags";
DROP TABLE "articles";