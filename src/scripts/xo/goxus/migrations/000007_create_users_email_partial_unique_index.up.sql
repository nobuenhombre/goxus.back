-- drop the old full unique index (if any) and recreate — this is idempotent
DROP INDEX IF EXISTS public.users_email_uindex;

-- recreate full unique index on email (code already prevents duplicates among all records)
CREATE UNIQUE INDEX users_email_uindex
    ON public.users (email);
