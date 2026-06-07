-- Rollback seed of 4 additional settings groups

-- 1. Delete settings for new groups
DELETE
FROM public.settings
WHERE group_id IN (SELECT id FROM public.settings_groups WHERE name IN ('Notifications', 'Privacy', 'Language & Region', 'Advanced'));

-- 2. Delete the groups themselves
DELETE
FROM public.settings_groups
WHERE name IN ('Notifications', 'Privacy', 'Language & Region', 'Advanced');