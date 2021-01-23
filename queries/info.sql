SELECT
  COUNT(1) AS blocks_count,
  MIN(height) AS blocks_min_height,
  MIN(timestamp) AS blocks_min_timestamp,
  MAX(height) AS blocks_max_height,
  MAX(timestamp) AS blocks_max_timestamp,
  COUNT(DISTINCT creator_id) AS blocks_producers_count,
  (SELECT COUNT(1) FROM public_keys) AS public_keys_count,
  (SELECT COUNT(1) FROM internal_commands) AS internal_commands_count,
  (SELECT COUNT(1) FROM user_commands) AS user_commands_count,
  {{ array }}
    SELECT type, COUNT(1)
    FROM user_commands
    GROUP BY type
  {{ end_array }} AS user_commands_types,
  {{ array }}
    SELECT type, COUNT(1)
    FROM internal_commands
    GROUP BY type
  {{ end_array }} AS internal_commands_types
FROM
  blocks