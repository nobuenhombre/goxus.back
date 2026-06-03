create table public.rbac_roles
(
    id         bigint generated always as identity
        constraint rbac_roles_pk
            primary key,
    name       varchar(255)            not null,
    slug       varchar(255)            not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

create unique index rbac_roles_slug_uindex
    on public.rbac_roles (slug);

create table public.rbac_permissions
(
    id         bigint generated always as identity
        constraint rbac_permissions_pk
            primary key,
    name       varchar(255)            not null,
    slug       varchar(255)            not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

create unique index rbac_permissions_slug_uindex
    on public.rbac_permissions (slug);

create table public.rbac_role_permissions
(
    id            bigint generated always as identity
        constraint rbac_role_permissions_pk
            primary key,
    role_id       bigint                  not null,
    permission_id bigint                  not null,
    created_at    timestamp default now() not null,
    updated_at    timestamp default now() not null
);

create unique index rbac_role_permissions_unique
    on public.rbac_role_permissions (role_id, permission_id);

alter table public.rbac_role_permissions
    add constraint rbac_role_permissions_roles_fk
        foreign key (role_id) references public.rbac_roles (id)
            on delete cascade;

alter table public.rbac_role_permissions
    add constraint rbac_role_permissions_permissions_fk
        foreign key (permission_id) references public.rbac_permissions (id)
            on delete cascade;

create table public.rbac_user_roles
(
    id         bigint generated always as identity
        constraint rbac_user_roles_pk
            primary key,
    user_id    bigint                  not null,
    role_id    bigint                  not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

create unique index rbac_user_roles_unique
    on public.rbac_user_roles (user_id, role_id);

alter table public.rbac_user_roles
    add constraint rbac_user_roles_roles_fk
        foreign key (role_id) references public.rbac_roles (id)
            on delete cascade;

alter table public.rbac_user_roles
    add constraint rbac_user_roles_users_fk
        foreign key (user_id) references public.users (id)
            on delete cascade;
