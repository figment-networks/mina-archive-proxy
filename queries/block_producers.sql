WITH transformed AS (
  SELECT
    json_build_object(
      public_keys.value,
      json_build_object(
        'blocks_produced', blocks_produced,
        'first_block', first_block,
        'last_block', last_block
      )
    ) AS obj
  FROM
    (
      SELECT
        creator_id,
        COUNT(1) AS blocks_produced,
        MIN(height) AS first_block,
        MAX(height) AS last_block
      FROM
        blocks
      GROUP BY
        creator_id
    ) all_creators
  INNER JOIN public_keys
    ON id = all_creators.creator_id
  ORDER BY
    blocks_produced DESC
)
SELECT
  json_object_agg(key, value)
FROM
  transformed, json_each(obj)