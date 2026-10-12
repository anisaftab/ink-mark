-- +goose Up
create table documents (
    id         text     primary key,
    title      text     not null default 'Untitled',
    file_path  text     not null,
    created_at datetime not null default (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at datetime not null default (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

create index documents_updated_at_idx on documents(updated_at desc);

-- +goose StatementBegin
create trigger documents_set_updated_at
    after update on documents
    for each row
    when old.title != new.title or old.file_path != new.file_path
begin
    update documents set updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
    where id = new.id;
end;
-- +goose StatementEnd

-- +goose Down
drop trigger if exists documents_set_updated_at;
drop table if exists documents;
