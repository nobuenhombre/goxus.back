create table public.users
(
    id                       bigint generated always as identity
        constraint users_pk
            primary key,
    name                     varchar(255)            not null,
    email                    varchar(255)            not null,
    password                 varchar(255)            not null,
    email_verified_at        timestamp,
    created_at               timestamp default now() not null,
    updated_at               timestamp default now() not null,
    deleted_at               timestamp
);

create unique index users_email_uindex
    on public.users (email);