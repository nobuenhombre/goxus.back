create table public.users_tokens
(
    id           bigint generated always as identity,
    token        varchar(255)            not null,
    user_id      bigint                  not null,
    last_used_at timestamp,
    created_at   timestamp default now() not null,
    updated_at   timestamp default now() not null,
    deleted_at   timestamp
);

-- primary key
alter table public.users_tokens
    add constraint users_tokens_pk primary key (id);

-- foreign key to users
alter table public.users_tokens
    add constraint users_tokens_users_fk foreign key (user_id) references public.users (id) on delete cascade;

-- unique token for fast lookup
create unique index users_tokens_token_uindex
    on public.users_tokens (token);

-- index for user token queries
create index users_tokens_user_id_index
    on public.users_tokens (user_id);