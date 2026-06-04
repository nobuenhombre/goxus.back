SELECT
    id, name, email, password, email_verified_at, created_at, updated_at, deleted_at
FROM
    public.users
WHERE
    email = %%email string%%
  AND deleted_at IS NULL