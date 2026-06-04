SELECT
    id, user_id, role_id, created_at, updated_at
FROM
    public.rbac_user_roles
WHERE
    role_id = %%roleID int64%%
ORDER BY
    id ASC