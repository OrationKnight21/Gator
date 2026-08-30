-- +goose Up
CREATE TABLE feed_follows(
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id uuid not null REFERENCES users(id) ON DELETE CASCADE,
    feed_id uuid not null REFERENCES feeds(id) ON DELETE CASCADE,
    unique(user_id,feed_id)
);
-- +goose Down
drop table feed_follows;