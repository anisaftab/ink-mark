-- +goose Up
alter table documents add column favorited integer not null default 0;

create index documents_favorited_idx on documents(favorited) where favorited = 1;

-- +goose Down
drop index if exists documents_favorited_idx;
