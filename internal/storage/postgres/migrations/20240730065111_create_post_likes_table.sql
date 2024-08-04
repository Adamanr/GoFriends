-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS post_likes (
    user_id int not NULL,
    post_id int references posts(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS post_likes;
-- +goose StatementEnd
