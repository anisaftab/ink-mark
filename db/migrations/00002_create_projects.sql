-- +goose Up
create table projects (
    id          text     primary key,
    name        text     not null,
    folder_path text     not null,
    created_at  datetime not null default (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at  datetime not null default (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

create index projects_updated_at_idx on projects(updated_at desc);

-- +goose StatementBegin
create trigger projects_set_updated_at
    after update on projects
    for each row
    when old.name != new.name or old.folder_path != new.folder_path
begin
    update projects set updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
    where id = new.id;
end;
-- +goose StatementEnd

alter table documents add column project_id text references projects(id) on delete set null;

create index documents_project_id_idx on documents(project_id);

-- +goose Down
drop index if exists documents_project_id_idx;
drop trigger if exists projects_set_updated_at;
drop index if exists projects_updated_at_idx;
drop table if exists projects;
