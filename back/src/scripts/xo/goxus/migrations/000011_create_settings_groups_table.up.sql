create table public.settings_groups
(
    id          bigint generated always as identity
        constraint settings_groups_pk
            primary key,
    name        varchar(255) not null,
    description text,
    "order"     bigint       not null default 0
);

comment on table public.settings_groups is 'Groups of user settings for UI organization';

-- Seed default appearance group
insert into public.settings_groups (name, description, "order")
values ('Appearance', 'Customize the appearance of the app. Automatically switch between day and night themes.', 1)
on conflict do nothing;