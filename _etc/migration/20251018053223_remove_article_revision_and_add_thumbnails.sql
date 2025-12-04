-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- WARNING: previous data won't be saved!
DROP TABLE "revisions";

ALTER TABLE "series"
    ADD COLUMN "thumbnail" VARCHAR(64) NOT NULL DEFAULT '';

ALTER TABLE "articles"
    ADD COLUMN "thumbnail" VARCHAR(64) NOT NULL DEFAULT '';

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

ALTER TABLE "articles"
    DROP COLUMN "thumbnail";

ALTER TABLE "series"
    DROP COLUMN "thumbnail";

CREATE TABLE "revisions"(
    "id" SERIAL PRIMARY KEY,
    "article_id" INTEGER NOT NULL,
    "message" VARCHAR(256) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY("article_id")
        REFERENCES "articles"("id")
        ON DELETE CASCADE);