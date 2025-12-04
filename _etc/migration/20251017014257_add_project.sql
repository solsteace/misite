-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE "projects"(
    "id" SERIAL PRIMARY KEY,
    "devblog_serie" INTEGER,

    "name" VARCHAR(128) NOT NULL,
    "thumbnail" VARCHAR(64) NOT NULL,
    "synopsis" VARCHAR(256) NOT NULL,
    "description" TEXT NOT NULL,

    FOREIGN KEY("devblog_serie")
        REFERENCES "series"("id")
        ON DELETE SET NULL);

CREATE TABLE "project_tags"(
    "id" SERIAL PRIMARY KEY,
    "tag_id" INTEGER NOT NULL,
    "project_id" INTEGER NOT NULL,

    FOREIGN KEY("tag_id")
        REFERENCES "tags"("id"),

    UNIQUE("tag_id", "project_id"));

CREATE TABLE "project_links"(
    "id" SERIAL PRIMARY KEY,
    "project_id" INTEGER NOT NULL,
    "display_text" VARCHAR(64) NOT NULL,
    "url" VARCHAR(128) NOT NULL,

    FOREIGN KEY("project_id")
        REFERENCES "projects"("id"),

    UNIQUE("project_id", "url"));

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE "project_links";
DROP TABLE "project_tags";
DROP TABLE "projects";