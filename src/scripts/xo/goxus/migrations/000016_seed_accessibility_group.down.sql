-- Rollback: delete Accessibility group and its settings

-- 1. Delete settings for the group
DELETE
FROM public.settings
WHERE group_id IN (SELECT id FROM public.settings_groups WHERE name = 'Accessibility');

-- 2. Delete the group itself
DELETE
FROM public.settings_groups
WHERE name = 'Accessibility';
