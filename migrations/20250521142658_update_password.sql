-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ALTER COLUMN password TYPE BYTEA USING password::BYTEA;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
