-- name: StoreCards :exec
INSERT INTO player_card_pool (player_id, name, set_code, collector_number, count)
SELECT unnest(@player_id::text[]), -- id
       unnest(@name::text[]), -- name
       unnest(@set_code::text[]), -- set_code
       unnest(@collector_number::int[]),  -- collector_number
       unnest(@count::int[])   -- count
ON CONFLICT (player_id, set_code, collector_number)
    DO UPDATE SET count = EXCLUDED.count + player_card_pool.count;

-- name: GetCards :many
SELECT name, set_code, collector_number, count
FROM player_card_pool
WHERE player_id = $1;

-- name: GetAllPlayers :many
SELECT id, wild_card_count, wild_pack_count
FROM player;

-- name: GetPlayer :one
SELECT id, wild_card_count, wild_pack_count
FROM player
WHERE id = $1;

-- name: UpdatePlayer :exec
UPDATE player
SET wild_card_count = $2,
    wild_pack_count = $3
WHERE ID = $1;

-- name: GetPairing :one
SELECT round, player_id1, player_id2, wins1, wins2, draws
from pairing
WHERE round = $1
  AND (player_id1 = $2 or player_id2 = $2);

-- name: StorePairings :exec
INSERT INTO pairing (round, player_id1, player_id2, wins1, wins2, draws)
SELECT unnest(@round::int[]),  -- round
       unnest(@player_id1::text[]), -- player_id1
       unnest(@player_id2::text[]), -- player_id2
       unnest(@wins1::int[]),  -- wins1
       unnest(@wins2::int[]),  -- wins2
       unnest(@draws::int[]);   -- draws
-- draws

-- name: UpdatePairing :exec
UPDATE pairing
SET wins1 = $4,
    wins2 = $5,
    draws = $6
WHERE round = $1
  AND player_id1 = $2
  AND player_id2 = $3
  AND wins1 = 0
  AND wins2 = 0
  AND draws = 0;

-- name: StartLeague :exec
INSERT INTO league (round, active, created_at)
VALUES (0, true, now());

-- name: EndLeague :exec
UPDATE league
SET active= false
WHERE active = true;

-- name: GetCurrentRound :one
SELECT round
FROM league
where active = true;