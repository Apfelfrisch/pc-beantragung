-- name: CreateSignOn :exec
INSERT INTO signons (
    id_pc, company, firstname, lastname, zip, city, street, house_no, pc_state, desired_delivery_start, meter_no, malo, melo, config_id, created_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetSignOn :one
SELECT sqlc.embed(signons), sqlc.embed(signon_context) FROM signons
LEFT JOIN signon_context ON signons.id_pc = signon_context.signon_id_pc
WHERE id = ? LIMIT 1;

-- name: ListSignOns :many
SELECT signons.* FROM signons
LEFT JOIN signon_context ON signons.id_pc = signon_context.signon_id_pc;

-- name: ListSignOnsByState :many
SELECT signons.* FROM signons
LEFT JOIN signon_context ON signons.id_pc = signon_context.signon_id_pc
WHERE signon_context.state = ?;

-- name: UpdateContext :exec
UPDATE signon_context SET state = ?, comment = ? WHERE signon_id_pc = ?;

-- name: FillupContext :exec
INSERT INTO signon_context(signon_id_pc, state)
select id_pc, "processing" from signons
where (
    select signon_id_pc from signon_context where signon_id_pc = id_pc
) is null;
