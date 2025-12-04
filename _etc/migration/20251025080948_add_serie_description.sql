-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE "series"
    ADD COLUMN "description" VARCHAR(256) NOT NULL DEFAULT '';

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

ALTER TABLE "series"
    DROP COLUMN "description";