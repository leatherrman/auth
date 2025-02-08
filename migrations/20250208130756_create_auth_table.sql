-- +goose Up
create table auth (
    id serial primary key,
    name text not null,
    email text not null,
    role INT not null,
    password text not null,
    created_at timestamp not null default now(),
    updated_at timestamp
);

-- +goose Down
drop table auth;
