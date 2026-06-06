-- create a view that extends users with concatenated role names
create or replace view public.user_with_roles as
select
    u.id,
    u.name,
    u.email,
    u.password,
    u.email_verified_at,
    u.created_at,
    u.updated_at,
    u.deleted_at,
    coalesce(
        (select string_agg(r.name, ', ' order by r.name)
         from public.rbac_roles r
         join public.rbac_user_roles ur on r.id = ur.role_id
         where ur.user_id = u.id),
        ''
    ) as roles
from public.users u
order by u.id asc;