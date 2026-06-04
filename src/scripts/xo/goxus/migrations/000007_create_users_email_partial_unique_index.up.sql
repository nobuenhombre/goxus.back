-- drop the old full unique index (if any) and recreate — this is idempotent
drop index if exists public.users_email_uindex;

-- recreate full unique index on email (code already prevents duplicates among all records)
create unique index users_email_uindex
    on public.users (email);