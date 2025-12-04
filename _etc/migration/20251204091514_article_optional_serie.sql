-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE "articles"
    ALTER COLUMN "serie_id" DROP NOT NULL;
ALTER TABLE "articles"
    ALTER COLUMN "serie_order" DROP NOT NULL;
ALTER TABLE "articles"
    ADD CONSTRAINT "serie_should_together_be_empty_or_null"
    CHECK ( 
        (serie_id IS NULL AND serie_order IS NULL)
        OR (serie_id IS NOT NULL AND serie_order IS NOT NULL));



-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

ALTER TABLE "articles"
    ALTER COLUMN "serie_id" ADD NOT NULL;
ALTER TABLE "articles"
    ALTER COLUMN "serie_order" ADD NOT NULL;
ALTER TABLE "articles"
    DROP CONSTRAINT "serie_should_together_be_empty_or_null";