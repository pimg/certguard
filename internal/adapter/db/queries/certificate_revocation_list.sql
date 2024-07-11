-- name: CreateCertificateRevocationList :one
INSERT INTO certificate_revocation_list(
    name,
    signature,
    this_update,
    next_update,
    url,
    raw
) VALUES (?,?,?,?,?,?)
  ON CONFLICT DO UPDATE SET
    this_update = excluded.this_update,
    next_update = excluded.next_update
RETURNING id;

-- name: UpdateCertificateRevocationList :one
UPDATE certificate_revocation_list
set this_update = ?,
    next_update = ?,
    raw = ?
WHERE name = ?
RETURNING *;

-- name: GetCertificateRevocationList :one
SELECT id, name, signature, DATETIME(this_update) as this_update, DATETIME(next_update) as next_update, url, raw FROM certificate_revocation_list
WHERE name = ?;

-- name: ListCertificateRevocationLists :many
SELECT id, name, signature, DATETIME(this_update) as this_update, DATETIME(next_update) as next_update, url, raw  FROM certificate_revocation_list
ORDER BY id;