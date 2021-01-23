WITH all_creators AS (
  SELECT creator_id, COUNT(1) AS blocks_produced
  FROM blocks
  GROUP BY creator_id
)
SELECT
  public_keys.value AS public_key, all_creators.blocks_produced
FROM
  all_creators
INNER JOIN public_keys
  ON id = all_creators.creator_id
ORDER BY
  blocks_produced DESC