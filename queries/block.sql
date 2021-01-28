SELECT
  blocks.height,
  blocks.state_hash,
  blocks.parent_hash,
  blocks.ledger_hash,
  snarked_ledger_hashes.value AS snarked_ledger_hash,
  creator_keys.value AS creator,
  winnder_keys.value AS winner,
  blocks.timestamp,
  to_char(to_timestamp(blocks.timestamp / 1000), 'YYYY-MM-DD"T"HH24:MI:SS:MS"Z"') AS timestamp_formatted,
  blocks.global_slot_since_genesis,
  blocks.global_slot,
  {{ array }}
    SELECT
      (hash || '-' || sequence_no || '-' || secondary_sequence_no) AS id,
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
      blocks_internal_commands.block_id = blocks.id
    ORDER BY
      blocks_internal_commands.secondary_sequence_no ASC, blocks_internal_commands.sequence_no ASC
  {{ end_array }} AS internal_commands,
  {{ array }}
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
      blocks_user_commands.block_id = blocks.id
    ORDER BY
      blocks_user_commands.user_command_id ASC
  {{ end_array }} AS user_commands
FROM
  blocks
INNER JOIN public_keys creator_keys
  ON creator_keys.id = blocks.creator_id
INNER JOIN public_keys winnder_keys
  ON winnder_keys.id = blocks.block_winner_id
INNER JOIN snarked_ledger_hashes
  ON snarked_ledger_hashes.id = blocks.snarked_ledger_hash_id
WHERE
  blocks.state_hash = $1