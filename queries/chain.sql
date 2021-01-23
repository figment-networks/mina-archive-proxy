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

SELECT
  height,
  state_hash,
  parent_hash,
  timestamp
FROM
  chain
WHERE
  height >= $1
ORDER BY
  height ASC
LIMIT
  $2