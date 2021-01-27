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
  (
    SELECT json_object_agg(key, value)
    FROM (
      SELECT json_build_object(type, COUNT(1)) AS data FROM user_commands GROUP BY type
    ) t, json_each(data)
  ) AS user_commands_types,
  (
    SELECT json_object_agg(key, value)
    FROM (
      SELECT json_build_object(type, COUNT(1)) AS data FROM internal_commands GROUP BY type
    ) t, json_each(data)
  ) AS internal_commands_types
FROM
  blocks