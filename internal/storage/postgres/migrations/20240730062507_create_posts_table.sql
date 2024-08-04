-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(250),
    body text,
    author_id int NOT NULL,
    images_name TEXT[],
    likes int,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS posts_user_id_idx ON posts(user_id);
CREATE INDEX IF NOT EXISTS posts_title_idx ON posts(title);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;
DROP INDEX IF EXISTS posts_user_id_idx;
DROP INDEX IF EXISTS posts_title_idx;
-- +goose StatementEnd
