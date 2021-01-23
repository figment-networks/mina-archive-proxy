SELECT
  hash,
  type,
  fee,
  token,
  receivers.value AS receiver,
  blocks_internal_commands.sequence_no,
  blocks_internal_commands.secondary_sequence_no
FROM
  blocks_internal_commands
INNER JOIN internal_commands
  ON internal_commands.id = blocks_internal_commands.internal_command_id
INNER JOIN public_keys receivers
  ON receivers.id = internal_commands.receiver_id
WHERE
  blocks_internal_commands.block_id = (
    SELECT id FROM blocks WHERE state_hash = ? LIMIT 1
  )
ORDER BY
  blocks_internal_commands.internal_command_id ASC