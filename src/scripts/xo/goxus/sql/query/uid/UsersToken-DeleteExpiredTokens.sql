UPDATE public.users_tokens
SET deleted_at = (NOW() AT TIME ZONE 'UTC'), updated_at = (NOW() AT TIME ZONE 'UTC')
WHERE deleted_at IS NULL
  AND (
          CASE
              WHEN last_used_at IS NOT NULL THEN last_used_at
              ELSE created_at
              END
          ) < (NOW() AT TIME ZONE 'UTC') - make_interval(days => %%ttlDays int%%)
