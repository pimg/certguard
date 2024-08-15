-- name: CreateRevokedCertificates :exec
INSERT INTO revoked_certificate(
    serialnumber,
    revocation_date,
    reason,
    revocation_list
) VALUES (
          ?,?,?,?
)
ON CONFLICT DO NOTHING;

-- name: GetRevokedCertificatesByRevocationList :many
SELECT id, serialnumber, DATETIME(revocation_date) as revocation_date, reason, revocation_list
FROM revoked_certificate
WHERE revocation_list = ?
ORDER BY revocation_date;

-- name: GetRevokedCertificate :one
SELECT cert.id, cert.serialnumber, cert.reason, DATETIME(cert.revocation_date) as revocation_date, crl.name AS revoked_by
FROM revoked_certificate as cert
JOIN certificate_revocation_list AS crl ON crl.id = cert.revocation_list
WHERE serialnumber = ?;
