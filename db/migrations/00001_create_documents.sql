-- +goose Up
create extension if not exists "pgcrypto";

create table documents (
    id          uuid        primary key default gen_random_uuid(),
    title       text        not null default 'Untitled',
    content     jsonb       not null default '{}',
    content_fts tsvector    generated always as (
                                to_tsvector('english', coalesce(title, ''))
                            ) stored,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now()
);

create index documents_content_fts_idx on documents using gin(content_fts);
create index documents_updated_at_idx  on documents(updated_at desc);

-- +goose StatementBegin
create or replace function set_updated_at()
returns trigger language plpgsql as $$
begin
    new.updated_at = now();
    return new;
end;
$$;
-- +goose StatementEnd

create trigger documents_set_updated_at
    before update on documents
    for each row execute function set_updated_at();

-- +goose Down
drop trigger if exists documents_set_updated_at on documents;
drop function if exists set_updated_at();
drop table if exists documents;
