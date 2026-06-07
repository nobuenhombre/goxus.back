-- Seed 1 additional settings group "Accessibility" with 5 settings using
-- previously unused settings_types
--
-- Unused types (not referenced by any existing settings):
--   inputTextField, inputPasswordField, inputFloatNumberField,
--   textareaField, inputIntSliderRange
--
-- We use all 5 = within the requested 3–7 range.

-- 1. Seed the new group
INSERT INTO public.settings_groups (name, description, order_pos)
VALUES ('Accessibility', 'Customize accessibility features for an improved experience.', 6)
ON CONFLICT DO NOTHING;

-- 2. Seed settings — Accessibility group (5 settings, one per unused type)

-- inputTextField — single-line text input
INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Screen Reader Announcements',
       'Custom text that screen readers announce when the page loads.',
       '{}'::JSON,
       '{
         "value": ""
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'inputTextField'
  AND sg.name = 'Accessibility'
ON CONFLICT DO NOTHING;

-- inputPasswordField — password input (stored encrypted in user_settings)
INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Voice Access Passphrase',
       'Passphrase required to enable voice control commands.',
       '{}'::JSON,
       '{
         "value": ""
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'inputPasswordField'
  AND sg.name = 'Accessibility'
ON CONFLICT DO NOTHING;

-- inputFloatNumberField — floating point number input
INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Font Size Scale',
       'Global font size multiplier (0.5 – 2.0) for better readability.',
       '{
         "value": {
           "min": 0.5,
           "max": 2.0,
           "step": 0.1
         }
       }'::JSON,
       '{
         "value": 1.0
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'inputFloatNumberField'
  AND sg.name = 'Accessibility'
ON CONFLICT DO NOTHING;

-- textareaField — multi-line text area
INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Custom CSS Overrides',
       'Additional CSS rules injected for custom styling (one rule per line).',
       '{}'::JSON,
       '{
         "value": ""
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'textareaField'
  AND sg.name = 'Accessibility'
ON CONFLICT DO NOTHING;

-- inputIntSliderRange — integer slider with range (start/end)
INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id, sg.id, 'Reading Speed Range',
       'Minimum and maximum words-per-minute for text-to-speech playback.',
       '{
         "value": {
           "min": 100,
           "max": 400,
           "step": 10
         }
       }'::JSON,
       '{
         "value": {
           "start": 150,
           "end": 300
         }
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'inputIntSliderRange'
  AND sg.name = 'Accessibility'
ON CONFLICT DO NOTHING;