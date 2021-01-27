SELECT
  value AS id,
  (SELECT COUNT(1) FROM user_commands WHERE source_id = public_keys.id) AS commands_sent,
  (SELECT COUNT(1) FROM user_commands WHERE receiver_id = public_keys.id) AS commands_received
FROM
  public_keys
WHERE
  value = ?