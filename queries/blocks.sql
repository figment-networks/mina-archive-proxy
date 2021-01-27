{{ if .Canonical }}
  WITH RECURSIVE chain AS (
  (
    SELECT id, parent_id, height, state_hash, parent_hash, timestamp
    FROM blocks b
    WHERE height = (SELECT MAX(height) FROM blocks)
    ORDER BY timestamp ASC
    LIMIT 1
  )
  UNION ALL
    (
      SELECT b.id, b.parent_id, b.height, b.state_hash, b.parent_hash, b.timestamp
      FROM blocks b
      INNER JOIN chain
        ON b.id = chain.parent_id
        AND chain.id <> chain.parent_id
    )
  )
{{ end }}

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
  blocks.global_slot
FROM
  blocks

{{ if .Canonical }}
  INNER JOIN chain
    ON chain.id = blocks.id
{{ end }}

INNER JOIN public_keys creator_keys
  ON creator_keys.id = blocks.creator_id
INNER JOIN public_keys winnder_keys
  ON winnder_keys.id = blocks.block_winner_id
INNER JOIN snarked_ledger_hashes
  ON snarked_ledger_hashes.id = blocks.snarked_ledger_hash_id
WHERE
  blocks.height >= $1
  {{ if .Creator }}
    AND creator_keys.value = '{{ .Creator }}'
  {{ end }}
ORDER BY
  blocks.height ASC
LIMIT
  $2