-- Seed 4 additional settings groups with random settings

-- 1. Seed groups
INSERT INTO public.settings_groups (name, description, order_pos)
VALUES ('Notifications', 'Configure how and when you receive notifications.', 2),
       ('Privacy', 'Manage your privacy and security preferences.', 3),
       ('Language & Region', 'Set language, timezone, and regional preferences.', 4),
       ('Advanced', 'Advanced and developer-oriented settings.', 5)
ON CONFLICT DO NOTHING;

-- 2. Seed settings — Notifications group (4 settings)
INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Email Notifications',
       'Choose how often you receive email notifications.',
       '{
         "value": {
           "1": "Immediately",
           "2": "Daily digest",
           "3": "Weekly digest",
           "4": "Never"
         }
       }'::JSON,
       '{
         "value": 2
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'listRadios'
  AND sg.name = 'Notifications'
ON CONFLICT DO NOTHING;

INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Push Notifications',
       'Receive push notifications in your browser.',
       '{
         "value": true
       }'::JSON,
       '{
         "value": true
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'switch'
  AND sg.name = 'Notifications'
ON CONFLICT DO NOTHING;

INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Sound Alerts',
       'Play a sound when a notification arrives.',
       '{
         "value": true
       }'::JSON,
       '{
         "value": false
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'switch'
  AND sg.name = 'Notifications'
ON CONFLICT DO NOTHING;

INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Mute Period',
       'Suppress notifications during the specified time range.',
       '{
         "value": {
           "1": "Disabled",
           "2": "Night (22:00 – 08:00)",
           "3": "Custom"
         }
       }'::JSON,
       '{
         "value": 2
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'listRadios'
  AND sg.name = 'Notifications'
ON CONFLICT DO NOTHING;

-- 3. Seed settings — Privacy group (5 settings)
INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Show Online Status',
       'Let other users see when you are online.',
       '{
         "value": true
       }'::JSON,
       '{
         "value": true
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'switch'
  AND sg.name = 'Privacy'
ON CONFLICT DO NOTHING;

INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Profile Visibility',
       'Control who can see your profile information.',
       '{
         "value": {
           "1": "Everyone",
           "2": "Registered users only",
           "3": "Only me"
         }
       }'::JSON,
       '{
         "value": 2
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'listRadios'
  AND sg.name = 'Privacy'
ON CONFLICT DO NOTHING;

INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Data Retention Period',
       'How long to keep your activity data before automatic deletion.',
       '{
         "value": {
           "30": "30 days",
           "90": "90 days",
           "180": "180 days",
           "365": "1 year"
         }
       }'::JSON,
       '{
        "value": {
          "30": false,
          "90": true,
          "180": false,
          "365": false
        }
      }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'listChecks'
  AND sg.name = 'Privacy'
ON CONFLICT DO NOTHING;

INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Activity Logging',
       'Record your actions in the system for audit purposes.',
       '{
         "value": true
       }'::JSON,
       '{
         "value": true
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'switch'
  AND sg.name = 'Privacy'
ON CONFLICT DO NOTHING;

INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Two-Factor Authentication',
       'Enforce two-factor authentication for account access.',
       '{
         "value": {
           "1": "Disabled",
           "2": "Require for new devices",
           "3": "Always require"
         }
       }'::JSON,
       '{
         "value": 1
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'listRadios'
  AND sg.name = 'Privacy'
ON CONFLICT DO NOTHING;

-- 4. Seed settings — Language & Region group (3 settings)
INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Language',
       'Select the interface language.',
       '{
         "value": {
           "en": "English",
           "ru": "Русский",
           "es": "Español",
           "fr": "Français",
           "de": "Deutsch",
           "zh": "中文"
         }
       }'::JSON,
       '{
         "value": "en"
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'selectWithSearch'
  AND sg.name = 'Language & Region'
ON CONFLICT DO NOTHING;

INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Timezone',
       'Set your local timezone for correct time display.',
       '{
         "value": {
           "UTC": "UTC (Coordinated Universal Time)",
           "EST": "EST (Eastern Standard Time)",
           "CST": "CST (Central Standard Time)",
           "PST": "PST (Pacific Standard Time)",
           "MSK": "MSK (Moscow Standard Time)",
           "CET": "CET (Central European Time)"
         }
       }'::JSON,
       '{
         "value": "UTC"
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'selectWithSearch'
  AND sg.name = 'Language & Region'
ON CONFLICT DO NOTHING;

INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Auto-detect Language',
       'Automatically detect your language from browser settings.',
       '{
         "value": true
       }'::JSON,
       '{
         "value": true
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'switch'
  AND sg.name = 'Language & Region'
ON CONFLICT DO NOTHING;

-- 5. Seed settings — Advanced group (6 settings)
INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Developer Mode',
       'Enable developer tools and debug information in the UI.',
       '{
         "value": true
       }'::JSON,
       '{
         "value": false
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'switch'
  AND sg.name = 'Advanced'
ON CONFLICT DO NOTHING;

INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Results Per Page',
       'Number of items displayed per page in lists and tables.',
       '{}'::JSON,
       '{
         "value": 25
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'inputIntNumberField'
  AND sg.name = 'Advanced'
ON CONFLICT DO NOTHING;

INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Log Level',
       'Minimum log level to display in the developer console.',
       '{
         "value": {
           "debug": "Debug",
           "info": "Info",
           "warn": "Warning",
           "error": "Error"
         }
       }'::JSON,
       '{
         "value": "info"
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'selectSimple'
  AND sg.name = 'Advanced'
ON CONFLICT DO NOTHING;

INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Auto-save',
       'Automatically save changes without manual confirmation.',
       '{
         "value": true
       }'::JSON,
       '{
         "value": true
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'switch'
  AND sg.name = 'Advanced'
ON CONFLICT DO NOTHING;

INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Session Timeout',
       'Minutes of inactivity before automatic logout.',
       '{
         "value": {
           "min": 15,
           "max": 120,
           "step": 15
         }
       }'::JSON,
       '{
         "value": 30
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'inputIntSlider'
  AND sg.name = 'Advanced'
ON CONFLICT DO NOTHING;

INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Experimental Features',
       'Opt in to pre-release features that may be unstable.',
       '{
         "value": true
       }'::JSON,
       '{
         "value": false
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'switch'
  AND sg.name = 'Advanced'
ON CONFLICT DO NOTHING;