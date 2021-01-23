SELECT
  hash,
  type,
  fee_token,
  token,
  nonce,
  amount,
  fee,
  valid_until,
  memo,
  status,
  failure_reason,
  fee_payer_account_creation_fee_paid,
  receiver_account_creation_fee_paid,
  created_token,
  blocks_user_commands.sequence_no,
  fee_payers.value AS fee_payer,
  senders.value AS sender,
  receivers.value AS receiver,
  blocks_user_commands.sequence_no
FROM
  blocks_user_commands
INNER JOIN user_commands
  ON user_commands.id = blocks_user_commands.user_command_id
INNER JOIN public_keys senders
  ON senders.id = user_commands.source_id
INNER JOIN public_keys fee_payers
  ON fee_payers.id = user_commands.fee_payer_id
INNER JOIN public_keys receivers
  ON receivers.id = user_commands.receiver_id
WHERE
  blocks_user_commands.block_id = (
    SELECT id FROM blocks WHERE state_hash = ? LIMIT 1
  )
ORDER BY
  blocks_user_commands.user_command_id ASC