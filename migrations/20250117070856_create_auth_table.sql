-- +goose Up
create table auth
(
    uuid       uuid primary key,
    name       text      not null,
    email      text      not null,
    password   text      not null,
    role       integer   not null default 1,
    created_at timestamp not null default now(),
    updated_at timestamp,
    constraint check_role check (
        role in (0, 1, 2)
        )
);
-- +goose Down
drop table auth;
