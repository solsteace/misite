-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE "projects"
    ALTER COLUMN "thumbnail" DROP NOT NULL;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

ALTER TABLE "projects"
    ALTER COLUMN "thumbnail" ADD NOT NULL;
